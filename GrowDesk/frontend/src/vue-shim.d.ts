declare module 'vue' {
  export interface TypedRef<T> {
    value: T;
  }
  export function ref<T>(value: T): TypedRef<T>;
} 