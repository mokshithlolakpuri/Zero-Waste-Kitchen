export interface Recipe {
    id?: number;
    title: string;
    ingredients: string;
    instructions: string;
    prep_time: number;
    cook_time: number;
    servings: number;
    difficulty: string;
    cuisine: string;
    user_id?: number;
    created_at?: string;
  }