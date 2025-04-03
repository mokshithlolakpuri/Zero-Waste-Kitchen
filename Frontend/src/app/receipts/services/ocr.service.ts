// src/app/services/ocr.service.ts
import { Injectable } from '@angular/core';
import { createWorker, Worker } from 'tesseract.js';
import { BehaviorSubject } from 'rxjs';

interface TesseractProgress {
  status: string;
  progress: number;
  // Add other properties you might use
}

@Injectable({
  providedIn: 'root'
})
export class OcrService {
  private progress$ = new BehaviorSubject<number>(0);
  private status$ = new BehaviorSubject<string>('');

  constructor() { }

  async recognize(image: File): Promise<{ text: string; confidence: number }> {
    this.status$.next('Initializing Tesseract...');
    this.progress$.next(0);

    const worker = await createWorker();

    worker.setParameters({
      logger: (m: TesseractProgress) => {
        if (m.status === 'recognizing text') {
          this.progress$.next(m.progress);
          this.status$.next(`Processing: ${Math.round(m.progress * 100)}%`);
        } else {
          this.status$.next(m.status);
        }
      }
    });

    try {
      await worker.load();
      // Use the correct method names from the Tesseract.js v4 API
      await worker.load();
      await worker.reinitialize('eng');
      
      this.status$.next('Processing receipt...');
      const { data } = await worker.recognize(image);
      
      return {
        text: data.text,
        confidence: data.confidence
      };
    } finally {
      await worker.terminate();
      this.progress$.next(1);
      this.status$.next('Done');
    }
  }

  getProgress() {
    return this.progress$.asObservable();
  }

  getStatus() {
    return this.status$.asObservable();
  }
}