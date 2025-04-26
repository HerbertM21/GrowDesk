/**
 * Interfaz para representar una pregunta frecuente (FAQ)
 */
export interface FAQ {
  id?: number;
  question: string;
  answer: string;
  category: string;
  isPublished: boolean;
  createdAt?: string;
  updatedAt?: string;
}

/**
 * Interfaz para crear una nueva FAQ
 */
export interface FaqCreateData {
  question: string;
  answer: string;
  category: string;
  isPublished: boolean;
}

/**
 * Interfaz para actualizar una FAQ existente
 */
export interface FaqUpdateData extends Partial<FAQ> {
  id: number;
}