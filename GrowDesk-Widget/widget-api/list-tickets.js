// Script para listar los tickets guardados y sus detalles
const fs = require('fs');
const path = require('path');

// Directorio donde se almacenan los tickets
const ticketsDir = path.join(__dirname, 'data');

// Función para leer un archivo de ticket
function readTicketFile(filename) {
  const filePath = path.join(ticketsDir, filename);
  try {
    const data = fs.readFileSync(filePath, 'utf8');
    return JSON.parse(data);
  } catch (error) {
    console.error(`Error al leer el archivo ${filename}:`, error.message);
    return null;
  }
}

// Función principal
function listTickets() {
  console.log('Listando tickets disponibles:');
  console.log('----------------------------');
  
  try {
    // Verificar si el directorio existe
    if (!fs.existsSync(ticketsDir)) {
      console.log('No hay tickets disponibles (el directorio data no existe)');
      return;
    }
    
    // Leer archivos en el directorio
    const files = fs.readdirSync(ticketsDir);
    
    if (files.length === 0) {
      console.log('No hay tickets disponibles');
      return;
    }
    
    // Filtrar solo archivos JSON de tickets
    const ticketFiles = files.filter(file => file.startsWith('ticket_') && file.endsWith('.json'));
    
    if (ticketFiles.length === 0) {
      console.log('No hay tickets disponibles');
      return;
    }
    
    // Leer cada archivo y mostrar detalles
    ticketFiles.forEach((file, index) => {
      const ticket = readTicketFile(file);
      if (ticket) {
        console.log(`[${index + 1}] ID: ${ticket.ID}`);
        console.log(`    Creado por: ${ticket.UserName || 'Anónimo'} (${ticket.UserEmail || 'email desconocido'})`);
        console.log(`    Creado: ${new Date(ticket.CreatedAt).toLocaleString()}`);
        console.log(`    Mensajes: ${ticket.Messages ? ticket.Messages.length : 0}`);
        console.log(`    Último mensaje: ${ticket.Messages && ticket.Messages.length > 0 ? 
          ticket.Messages[ticket.Messages.length - 1].Content : 'N/A'}`);
        console.log('----------------------------');
      }
    });
    
    console.log(`Total de tickets: ${ticketFiles.length}`);
  } catch (error) {
    console.error('Error al listar tickets:', error.message);
  }
}

// Ejecutar la función principal
listTickets(); 