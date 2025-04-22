// Script para probar el envío de mensajes desde el API a los clientes WebSocket
const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args));
const fs = require('fs');
const path = require('path');

// Función para enviar mensaje
async function sendAgentMessage(ticketId, content, agentName = 'Soporte') {
  const message = {
    ticketId,
    content,
    agentName
  };

  console.log(`Enviando mensaje a ticket ${ticketId}: "${content}"`);

  try {
    const response = await fetch('http://localhost:8082/api/agent/messages', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Widget-ID': 'demo-widget',
        'X-Widget-Token': 'demo-token'
      },
      body: JSON.stringify(message)
    });

    // Verificar si la respuesta fue exitosa
    if (!response.ok) {
      const text = await response.text();
      console.error(`Error en la respuesta: ${response.status} ${response.statusText}`);
      console.error(`Respuesta: ${text}`);
      throw new Error(`Error del servidor: ${response.status}`);
    }

    // Intentar parsear como JSON, si falla mostrar el texto
    let data;
    try {
      data = await response.json();
    } catch (error) {
      const text = await response.text();
      console.log('Respuesta (no JSON):', text);
      return { success: true, response: text };
    }
    
    console.log('Respuesta del servidor:', data);
    return data;
  } catch (error) {
    console.error('Error al enviar mensaje:', error);
    throw error;
  }
}

// Función para crear un ticket de prueba si no existe ninguno
async function createTestTicketIfNeeded() {
  // Verificar si ya hay tickets
  const dataDir = path.join(__dirname, 'data');
  if (!fs.existsSync(dataDir)) {
    fs.mkdirSync(dataDir);
  }
  
  const files = fs.readdirSync(dataDir);
  const ticketFiles = files.filter(file => file.startsWith('ticket_') && file.endsWith('.json'));
  
  if (ticketFiles.length > 0) {
    // Ya hay tickets, obtener el primer ID
    const ticketFile = ticketFiles[0];
    const ticketData = JSON.parse(fs.readFileSync(path.join(dataDir, ticketFile), 'utf8'));
    console.log(`Usando ticket existente: ${ticketData.ID}`);
    return ticketData.ID;
  }
  
  // No hay tickets, crear uno nuevo
  console.log('No hay tickets, creando uno de prueba...');
  
  try {
    const response = await fetch('http://localhost:8082/widget/tickets', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-Widget-ID': 'demo-widget',
        'X-Widget-Token': 'demo-token'
      },
      body: JSON.stringify({
        name: 'Usuario de Prueba',
        email: 'test@example.com',
        message: 'Este es un ticket de prueba creado automáticamente.',
        metadata: {
          url: 'http://localhost:8090',
          referrer: '',
          userAgent: 'Test Script',
          screenSize: '1920x1080'
        }
      })
    });
    
    const data = await response.json();
    console.log('Ticket creado:', data);
    return data.ticketId;
  } catch (error) {
    console.error('Error al crear ticket de prueba:', error);
    // Usar un ID fijo en caso de error
    return 'TICKET-TEST123';
  }
}

// Función principal
async function main() {
  try {
    // Obtener ticketId de los argumentos, o crear uno si no existe
    let ticketId = process.argv[2];
    
    if (!ticketId) {
      ticketId = await createTestTicketIfNeeded();
    }
    
    // Obtener mensaje de los argumentos, o usar uno predeterminado
    const mensaje = process.argv[3] || '¡Hola! Soy un agente de soporte. ¿En qué puedo ayudarte?';
    
    console.log(`Usando ticket ID: ${ticketId}`);
    const result = await sendAgentMessage(ticketId, mensaje);
    console.log('Mensaje enviado con éxito:', result);
  } catch (error) {
    console.error('Error en el proceso:', error);
    process.exit(1);
  }
}

// Ejecutar
main(); 