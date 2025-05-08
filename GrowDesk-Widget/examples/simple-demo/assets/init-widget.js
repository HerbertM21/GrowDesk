// Inicialización del widget de soporte GrowDesk
(function() {
  // Log de inicialización
  console.log("Iniciando inicialización del widget");
  
  // Referencia para almacenar la instancia del widget
  let widgetInstance = null;
  
  // Configuración global
  window.GrowDeskConfig = {
    widgetId: 'demo-widget',
    apiUrl: 'http://localhost:3001',
    brandName: 'GrowDesk Demo',
    welcomeMessage: '¡Hola! ¿En qué podemos ayudarte hoy?',
    position: 'bottom-right'
  };
  
  console.log("Configuración del widget establecida", window.GrowDeskConfig);
  
  // Método temporal para abrir el widget antes de que se inicialice completamente
  window.GrowDeskWidget = {
    open: function() {
      console.log("Intento de abrir widget antes de inicialización completa");
      // Almacenamos la intención para ejecutarla cuando el widget esté listo
      window.GrowDeskWidget._tempOpen = true;
      console.log("Llamada temprana a open(), usando _tempOpen");
    }
  };
  
  // Verificar si ya existe el objeto GrowDeskWidget
  console.log("Objeto GrowDeskWidget pre-inicializado:", window.GrowDeskWidget);
})();
