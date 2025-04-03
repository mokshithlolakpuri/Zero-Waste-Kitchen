// Import the Firebase scripts required for messaging
importScripts('https://www.gstatic.com/firebasejs/9.23.0/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/9.23.0/firebase-messaging-compat.js');

// Initialize Firebase with your configuration
const firebaseConfig = {
  apiKey: "AIzaSyC60x-l8q4674Hb31HQ2AhTHiYHBzpeJ14",
  authDomain: "zero-kitchen-waste.firebaseapp.com",
  projectId: "zero-kitchen-waste",
  storageBucket: "zero-kitchen-waste.firebasestorage.app",
  messagingSenderId: "878484942907",
  appId: "1:878484942907:web:bddf19236a07e834b1ac1e",
  measurementId: "G-BW3GP4B7K9",
  vapidKey: "BMwMi4oOSDGUGu2JtXnCqSZagA7GmezicDHtRPcNsIkudssQXxovZB3HZWuga43Mgzd56d8f58pT66ZBZyczBGg"
};


firebase.initializeApp(firebaseConfig);

// Retrieve Firebase Messaging instance
const messaging = firebase.messaging();

// Handle background messages
// messaging.onBackgroundMessage((payload) => {
//   console.log('[firebase-messaging-sw.js] Received background message ', payload);

//   const notificationTitle = payload.notification?.title || 'Background Notification';
//   const notificationOptions = {
//     body: payload.notification?.body || 'You have a new message.',
//     icon: '/assets/icons/icon-192x192.png' // Replace with your app's icon
//   };

//   self.registration.showNotification(notificationTitle, notificationOptions);
// });