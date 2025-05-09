/**
 * GrowDesk Anti DevTools Script
 * Este script detecta y bloquea el uso de herramientas de desarrollo del navegador
 * excepto para administradores
 */
(function() {
  let devtoolsDetected = false;
  const originalTitle = document.title;
  const warningTitle = "⚠️ Acceso Bloqueado - GrowDesk";
  
  // Función para verificar si el usuario es administrador
  function isUserAdmin() {
    try {
      // Obtener datos del usuario desde localStorage
      const userId = localStorage.getItem('userId');
      const usersJson = localStorage.getItem('growdesk-users');
      
      if (!userId || !usersJson) {
        return false;
      }
      
      const users = JSON.parse(usersJson);
      
      // Encontrar el usuario actual
      const currentUser = users.find(user => user.id.toString() === userId.toString());
      
      // Verificar si el usuario es administrador
      return currentUser && currentUser.role === 'admin';
    } catch (e) {
      console.error('Error al verificar si el usuario es administrador:', e);
      return false;
    }
  }
  
  // Detectar DevTools por diferencia de tamaño
  function detectDevToolsBySize() {
    const widthThreshold = window.outerWidth - window.innerWidth > 160;
    const heightThreshold = window.outerHeight - window.innerHeight > 160;
    return widthThreshold || heightThreshold;
  }
  
  // Detectar DevTools por disponibilidad de funciones de depuración
  function detectDevToolsByDebugger() {
    let devtoolsOpen = false;
    
    // Firefox & Chrome
    const dateCheck = new Date();
    debugger;
    const diff = new Date() - dateCheck;
    if (diff > 100) {
      devtoolsOpen = true;
    }
    
    // Firebug
    if (window.Firebug && window.Firebug.chrome && window.Firebug.chrome.isInitialized) {
      devtoolsOpen = true;
    }
    
    // Chrome, Edge & otros basados en Chromium
    if ((window.console && window.console.firebug) || 
        (navigator.userAgent.toLowerCase().indexOf('firefox') > -1 && window.devtools && window.devtools.open)) {
      devtoolsOpen = true;
    }
    
    return devtoolsOpen;
  }
  
  // Detectar por tiempos de ejecución
  function detectDevToolsByTiming() {
    const start = performance.now();
    for (let i = 0; i < 1000; i++) {
      console.log(i);
      console.clear();
    }
    const end = performance.now();
    return (end - start) > 100; // Si tarda más de 100ms, probablemente las DevTools están abiertas
  }
  
  // Bloquear el contenido cuando se detectan DevTools
  function blockContent() {
    // Si el usuario es administrador, permitir el uso de DevTools
    if (isUserAdmin()) {
      // No bloqueamos para administradores, pero podemos mostrar un mensaje en consola
      console.info("%c✅ Modo desarrollador activado - Acceso permitido para administradores", "color: #5e35b1; font-weight: bold; font-size: 14px;");
      return;
    }
    
    if (devtoolsDetected) {
      return; // Ya estaba bloqueado
    }
    
    // Guardar el contenido original para restaurarlo si se cierran las DevTools
    const originalBody = document.body.innerHTML;
    
    // Cambiar título
    document.title = warningTitle;
    
    // Bloquear UI
    document.body.innerHTML = `
      <div id="devtools-blocker" style="
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100vh;
        background: linear-gradient(135deg, #5e35b1, #7b1fa2);
        color: white;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        z-index: 99999;
        font-family: 'Segoe UI', Arial, sans-serif;
        padding: 20px;
        text-align: center;
      ">
        <div style="
          background-color: rgba(255, 255, 255, 0.1);
          border-radius: 12px;
          padding: 40px;
          backdrop-filter: blur(5px);
          box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
          max-width: 500px;
        ">
          <svg xmlns="http://www.w3.org/2000/svg" width="64" height="64" viewBox="0 0 24 24" style="margin-bottom: 20px;">
            <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
          </svg>
          <h1 style="font-size: 28px; margin-bottom: 20px; color: white; text-shadow: 0 2px 4px rgba(0,0,0,0.3);">Acceso Denegado</h1>
          <p style="font-size: 18px; line-height: 1.6; margin-bottom: 25px; color: rgba(255, 255, 255, 0.9);">
            Por razones de seguridad, el uso de herramientas de desarrollo no está permitido en GrowDesk.
            Por favor, cierre las herramientas de desarrollo para continuar utilizando la aplicación.
          </p>
          <button id="restore-button" style="
            padding: 12px 24px;
            background-color: white;
            color: #5e35b1;
            border: none;
            border-radius: 50px;
            font-size: 16px;
            font-weight: bold;
            cursor: pointer;
            transition: all 0.3s ease;
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
          ">Continuar</button>
          
          <div style="margin-top: 20px; font-size: 14px; opacity: 0.8;">
            <p style="margin: 8px 0;">Si es administrador, puede iniciar sesión para acceder a las herramientas de desarrollo.</p>
          </div>
        </div>
      </div>
    `;
    
    // Agregar evento al botón para intentar restaurar
    setTimeout(() => {
      const button = document.getElementById('restore-button');
      if (button) {
        button.addEventListener('mouseover', function() {
          this.style.transform = 'translateY(-3px)';
          this.style.boxShadow = '0 6px 16px rgba(0, 0, 0, 0.2)';
        });
        
        button.addEventListener('mouseout', function() {
          this.style.transform = 'translateY(0)';
          this.style.boxShadow = '0 4px 12px rgba(0, 0, 0, 0.15)';
        });
        
        button.addEventListener('click', function() {
          // Verificar si las DevTools siguen abiertas
          if (!detectDevToolsBySize() && !detectDevToolsByDebugger()) {
            document.body.innerHTML = originalBody;
            document.title = originalTitle;
            devtoolsDetected = false;
            initProtection(); // Reiniciar la protección
          } else {
            alert('Por favor, cierre las herramientas de desarrollo primero.');
          }
        });
      }
    }, 100);
    
    devtoolsDetected = true;
  }
  
  // Restaurar contenido si se cierran las DevTools
  function checkAndRestore() {
    if (devtoolsDetected) {
      if (!detectDevToolsBySize() && !detectDevToolsByDebugger()) {
        location.reload(); // Recargar la página si se cierran las DevTools
      }
    }
  }
  
  // Deshabilitar teclas comunes para acceder a DevTools
  function disableDevToolsShortcuts() {
    document.addEventListener('keydown', function(e) {
      // Si el usuario es administrador, permitir las teclas de DevTools
      if (isUserAdmin()) {
        return true;
      }
      
      // Deshabilitar F12, Ctrl+Shift+I, Ctrl+Shift+J, Ctrl+Shift+C, Ctrl+U
      if (
        e.keyCode === 123 || // F12
        (e.ctrlKey && e.shiftKey && (e.keyCode === 73 || e.keyCode === 74 || e.keyCode === 67)) || // Ctrl+Shift+I/J/C
        (e.ctrlKey && e.keyCode === 85) // Ctrl+U
      ) {
        e.preventDefault();
        blockContent();
        return false;
      }
    });
    
    // Deshabilitar clic derecho
    document.addEventListener('contextmenu', function(e) {
      // Si el usuario es administrador, permitir clic derecho
      if (isUserAdmin()) {
        return true;
      }
      
      e.preventDefault();
      return false;
    });
  }
  
  // Iniciar protección
  function initProtection() {
    // Comprobar periódicamente si las DevTools están abiertas
    setInterval(function() {
      // No bloquear si el usuario es administrador
      if (isUserAdmin()) {
        return;
      }
      
      if (detectDevToolsBySize() || detectDevToolsByDebugger()) {
        blockContent();
      } else {
        checkAndRestore();
      }
    }, 1000);
    
    // Deshabilitar atajos de teclado para DevTools
    disableDevToolsShortcuts();
    
    // Sobreescribir console.log y otras funciones de consola solo para no-administradores
    if (!devtoolsDetected && !isUserAdmin()) {
      const noop = function() {};
      const methods = ['log', 'debug', 'info', 'warn', 'error', 'assert', 'dir', 'dirxml', 'trace', 'group', 'groupCollapsed', 'groupEnd', 'time', 'timeEnd', 'profile', 'profileEnd', 'count', 'clear', 'table', 'console', 'error'];
      
      methods.forEach(function(method) {
        console[method] = noop;
      });
    }
  }
  
  // Iniciar cuando el DOM esté listo
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initProtection);
  } else {
    initProtection();
  }
})(); 