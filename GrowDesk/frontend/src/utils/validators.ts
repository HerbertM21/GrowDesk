/**
 * Funciones de validación para objetos y datos
 */

import type { User } from '@/stores/users';
import type { Ticket, Tag } from '@/stores/tickets';

// Definir el tipo Category localmente para evitar problemas de importación
interface Category {
  id: number;
  name: string;
  description: string;
}

/**
 * Valida si un objeto cumple con los requisitos mínimos de un usuario
 */
export const isValidUser = (item: any): item is User => {
  return typeof item === 'object' && 
         item !== null &&
         'id' in item && 
         'email' in item && 
         'firstName' in item && 
         'lastName' in item && 
         'role' in item;
};

/**
 * Filtra un array para dejar solo usuarios válidos
 */
export const filterValidUsers = (items: any[]): User[] => {
  if (!Array.isArray(items)) return [];
  return items.filter(isValidUser);
};

/**
 * Valida si un objeto cumple con los requisitos mínimos de un ticket
 */
export const isValidTicket = (item: any): item is Ticket => {
  return typeof item === 'object' && 
         item !== null &&
         'id' in item &&
         'title' in item;
};

/**
 * Filtra un array para dejar solo tickets válidos
 */
export const filterValidTickets = (items: any[]): Ticket[] => {
  if (!Array.isArray(items)) return [];
  return items.filter(isValidTicket);
};

/**
 * Valida si un objeto cumple con los requisitos mínimos de una categoría
 */
export const isValidCategory = (item: any): item is Category => {
  return typeof item === 'object' && 
         item !== null &&
         'id' in item &&
         'name' in item;
};

/**
 * Filtra un array para dejar solo categorías válidas
 */
export const filterValidCategories = (items: any[]): Category[] => {
  if (!Array.isArray(items)) return [];
  return items.filter(isValidCategory);
};

/**
 * Limpia un objeto de localStorage y retorna solo los objetos válidos
 */
export const cleanLocalStorageData = <T>(key: string, validator: (item: any) => boolean): T[] => {
  try {
    const data = localStorage.getItem(key);
    if (!data) return [];
    
    const parsed = JSON.parse(data);
    if (!Array.isArray(parsed)) return [];
    
    const valid = parsed.filter(validator);
    
    // Si se filtraron elementos, actualizar localStorage
    if (valid.length !== parsed.length) {
      console.warn(`Se filtraron ${parsed.length - valid.length} elementos inválidos de ${key}`);
      localStorage.setItem(key, JSON.stringify(valid));
    }
    
    return valid as T[];
  } catch (err) {
    console.error(`Error al limpiar datos de ${key}:`, err);
    return [];
  }
}; 