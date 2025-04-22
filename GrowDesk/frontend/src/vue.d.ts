// Este archivo corrige el error "Cannot find module 'vue'"
declare module 'vue' {
  import { DefineComponent } from '@vue/runtime-core'
  
  // Define los componentes de Vue
  export const defineComponent: any
  export const ref: any
  export const reactive: any
  export const computed: any
  export const watch: any
  export const onMounted: any
  export const createApp: any
  
  // Define el tipo de componente
  export type Component = any
  
  // Exporta DefineComponent
  export { DefineComponent }
} 