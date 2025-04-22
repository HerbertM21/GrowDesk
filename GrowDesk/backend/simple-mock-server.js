const http = require('http');
const fs = require('fs');
const path = require('path');
const WebSocket = require('ws');

const PORT = process.env.PORT || 8000;
const TICKETS_FILE = path.join(__dirname, 'tickets.json');
const USERS_FILE = path.join(__dirname, 'users.json');
const CATEGORIES_FILE = path.join(__dirname, 'categories.json');

// Cargar tickets desde archivo o inicializar array vacío
let tickets = [];
try {
  if (fs.existsSync(TICKETS_FILE)) {
    const ticketsData = fs.readFileSync(TICKETS_FILE, 'utf8');
    tickets = JSON.parse(ticketsData);
    console.log(`Cargados ${tickets.length} tickets desde archivo`);
  } else {
    console.log('No se encontró archivo de tickets, iniciando con lista vacía');
  }
} catch (err) {
  console.error('Error al cargar tickets:', err);
}

// Cargar usuarios desde archivo o inicializar con datos de muestra
let users = [];
try {
  if (fs.existsSync(USERS_FILE)) {
    const usersData = fs.readFileSync(USERS_FILE, 'utf8');
    users = JSON.parse(usersData);
    console.log(`Cargados ${users.length} usuarios desde archivo`);
  } else {
    // Inicializar con usuarios de ejemplo
    users = [
      {
        id: '1',
        firstName: 'Admin',
        lastName: 'Usuario',
        email: 'admin@example.com',
        role: 'admin',
        department: 'Tecnología',
        active: true,
        password: 'password' // SE HASHEARA
      },
      {
        id: '2',
        firstName: 'Asistente',
        lastName: 'Soporte',
        email: 'asistente@example.com',
        role: 'assistant',
        department: 'Soporte',
        active: true,
        password: 'password'
      },
      {
        id: '3',
        firstName: 'Empleado',
        lastName: 'Regular',
        email: 'empleado@example.com',
        role: 'employee',
        department: 'Ventas',
        active: true,
        password: 'password'
      }
    ];
    console.log('Inicializando con usuarios de muestra');
  }
} catch (err) {
  console.error('Error al cargar usuarios:', err);
}

// Cargar categorías desde archivo o inicializar con datos de muestra
let categories = [];
try {
  if (fs.existsSync(CATEGORIES_FILE)) {
    const categoriesData = fs.readFileSync(CATEGORIES_FILE, 'utf8');
    categories = JSON.parse(categoriesData);
    console.log(`Cargadas ${categories.length} categorías desde archivo`);
  } else {
    // Inicializar con categorías de ejemplo
    categories = [
      {
        id: '1',
        name: 'Soporte Técnico',
        description: 'Problemas técnicos con hardware o software',
        color: '#4CAF50',
        icon: 'computer',
        active: true,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      },
      {
        id: '2',
        name: 'Consultas Generales',
        description: 'Preguntas sobre productos o servicios',
        color: '#2196F3',
        icon: 'help',
        active: true,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      },
      {
        id: '3',
        name: 'Facturación',
        description: 'Problemas o dudas sobre facturación',
        color: '#FFC107',
        icon: 'credit_card',
        active: true,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      }
    ];
    console.log('Inicializando con categorías de muestra');
  }
} catch (err) {
  console.error('Error al cargar categorías:', err);
}

// Función para guardar tickets en archivo
const saveTickets = () => {
  try {
    fs.writeFileSync(TICKETS_FILE, JSON.stringify(tickets, null, 2));
    console.log(`Guardados ${tickets.length} tickets en archivo`);
  } catch (err) {
    console.error('Error al guardar tickets:', err);
  }
};

// Definir función para guardar usuarios
const saveUsers = () => {
  try {
    fs.writeFileSync(USERS_FILE, JSON.stringify(users, null, 2));
    console.log(`Guardados ${users.length} usuarios en archivo`);
  } catch (err) {
    console.error('Error al guardar usuarios:', err);
  }
};

// Función para guardar categorías en archivo
const saveCategories = () => {
  try {
    fs.writeFileSync(CATEGORIES_FILE, JSON.stringify(categories, null, 2));
    console.log(`Guardadas ${categories.length} categorías en archivo`);
  } catch (err) {
    console.error('Error al guardar categorías:', err);
  }
};

// Intentar guardar datos iniciales si no existen
if (users.length > 0 && !fs.existsSync(USERS_FILE)) {
  saveUsers();
}

if (categories.length > 0 && !fs.existsSync(CATEGORIES_FILE)) {
  saveCategories();
}

// Token simple para pruebas
const generateToken = () => {
  return 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbi0xMjMiLCJlbWFpbCI6ImFkbWluQGdyb3dkZXNrLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDA4ODQwMH0.8J5ayPvA4B-1vF5NaqRXycW1pmIl9qjKP6Ddj4Ot_Cw';
};

// Extraer token de cabeceras
const extractToken = (req) => {
  const authHeader = req.headers['authorization'];
  if (authHeader && authHeader.startsWith('Bearer ')) {
    return authHeader.substring(7);
  }
  return null;
};

// Validar token (simple)
const validateToken = (token) => {
  return token && token.length > 0;
};

// Leer el cuerpo de una solicitud
const readBody = async (req) => {
  return new Promise((resolve, reject) => {
    let body = '';
    req.on('data', chunk => {
      body += chunk.toString();
    });
    
    req.on('end', () => {
      try {
        if (body) {
          resolve(JSON.parse(body));
        } else {
          resolve({});
        }
      } catch (e) {
        reject(e);
      }
    });
    
    req.on('error', (err) => {
      reject(err);
    });
  });
};

// Almacenar conexiones de WebSocket por ticket
const ticketConnections = {};
// Mapa de conexiones alternativas para compatibilidad
const alternateConnectionMap = {};

// Crear el servidor
const server = http.createServer(async (req, res) => {
  // Log de la URL y método para depuración
  console.log(`${req.method} ${req.url}`);
  
  // Configurar CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-Requested-With, X-Widget-ID, X-Widget-Token, Accept, Cache-Control, Pragma, Expires');
  res.setHeader('Access-Control-Allow-Credentials', 'true');
  
  // Manejar solicitudes OPTIONS preflight
  if (req.method === 'OPTIONS') {
    res.writeHead(204);
    res.end();
    return;
  }

  try {
    // Verificar estado del servidor
    if (req.method === 'GET' && req.url === '/api/health') {
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ status: 'ok', message: 'Server running' }));
      return;
    }

    // Login
    if (req.method === 'POST' && req.url === '/api/auth/login') {
      const data = await readBody(req);
      
      if (data.email && data.password) {
        const token = generateToken();
        
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({
          token: token,
          user: {
            id: 'admin-123',
            email: data.email,
            firstName: 'Admin',
            lastName: 'User',
            role: 'admin'
          }
        }));
      } else {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Email y contraseña requeridos' }));
      }
      return;
    }

    // Información del usuario autenticado
    if (req.method === 'GET' && req.url === '/api/auth/me') {
      const token = extractToken(req);
      
      if (validateToken(token)) {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({
          id: 'admin-123',
          email: 'admin@growdesk.com',
          firstName: 'Admin',
          lastName: 'User',
          role: 'admin'
        }));
      } else {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
      }
      return;
    }

    // Obtener tickets
    if (req.method === 'GET' && req.url === '/api/tickets') {
      const token = extractToken(req);
      
      if (validateToken(token)) {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(tickets));
      } else {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
      }
      return;
    }

    // Obtener ticket específico por ID
    if (req.method === 'GET' && req.url.startsWith('/api/tickets/')) {
      const token = extractToken(req);
      
      if (validateToken(token)) {
        // Comprobar si es una solicitud de mensajes
        if (req.url.includes('/messages')) {
          const ticketId = req.url.split('/api/tickets/')[1].split('/messages')[0];
          const ticket = tickets.find(t => t.id === ticketId);
          
          if (ticket) {
            res.writeHead(200, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify(ticket.messages || []));
          } else {
            res.writeHead(404, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'Ticket no encontrado' }));
          }
        } else {
          // Es una solicitud del ticket completo
          const ticketId = req.url.split('/api/tickets/')[1];
          const ticket = tickets.find(t => t.id === ticketId);
          
          if (ticket) {
            res.writeHead(200, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify(ticket));
          } else {
            res.writeHead(404, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'Ticket no encontrado' }));
          }
        }
      } else {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
      }
      return;
    }

    // Widget: crear ticket
    if (req.method === 'POST' && req.url === '/api/widget/tickets') {
      const data = await readBody(req);
      
      // Crear un nuevo ticket
      const newTicket = {
        id: 'TICKET-' + new Date().toISOString().replace(/[^0-9]/g, '').slice(0, 14),
        title: 'Soporte Web - ' + (data.name || 'Anónimo'),
        status: 'open',
        createdAt: new Date().toISOString(),
        customer: {
          name: data.name || 'Anónimo',
          email: data.email || 'anonymous@example.com'
        },
        messages: [
          {
            id: 'MSG-' + Date.now(),
            content: data.description || data.message || '',
            isClient: true,
            timestamp: new Date().toISOString()
          }
        ],
        description: data.description || data.message || '',
        priority: 'medium',
        category: 'soporte',
        createdBy: data.email || 'anonymous@example.com',
        assignedTo: null,
        updatedAt: new Date().toISOString()
      };
      
      // Añadir el ticket a la lista
      tickets.push(newTicket);
      
      // Guardar tickets en archivo
      saveTickets();
      
      res.writeHead(201, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({
        id: newTicket.id,
        message: 'Ticket creado con éxito',
        status: 'open'
      }));
      return;
    }

    // Widget: enviar mensaje a un ticket específico
    if (req.method === 'POST' && req.url.match(/^\/api\/tickets\/[A-Za-z0-9-]+\/messages$/)) {
      const token = extractToken(req);
      
      if (validateToken(token)) {
        const ticketId = req.url.split('/api/tickets/')[1].split('/messages')[0];
        const data = await readBody(req);
        const ticket = tickets.find(t => t.id === ticketId);
        
        if (ticket) {
          // FORZAR DEPURACIÓN DETALLADA
          console.log('================ DEBUG MENSAJE ================');
          console.log('URL:', req.url);
          console.log('Headers completos:', req.headers);
          console.log('Parámetros URL:', req.url.includes('from_client=true'));
          console.log('Parámetro from_client:', req.url.includes('from_client=true'));
          console.log('Header x-widget-id:', req.headers['x-widget-id']);
          console.log('Header x-message-source:', req.headers['x-message-source']);
          console.log('Contenido del mensaje:', data);
          
          // MEJORE LOGICA
          // Por defecto, mensajes desde el panel admin son de agente (isClient=false)
          // Mensajes con indicador de widget o from_client=true son del cliente (isClient=true)
          let isClientMessage = false;
          
          if (req.url.includes('from_client=true') || 
              req.headers['x-widget-id'] || 
              req.headers['x-message-source'] === 'widget-client') {
            console.log('Mensaje detectado como proveniente del WIDGET/CLIENTE');
            isClientMessage = true;
          } else {
            console.log('Mensaje detectado como proveniente del PANEL/AGENTE');
            isClientMessage = false;
          }
          
          console.log('DECISIÓN FINAL - isClient:', isClientMessage);
          
          // Crear nuevo mensaje con la bandera isClient determinada
          const newMessage = {
            id: 'MSG-' + Date.now(),
            content: data.content,
            isClient: isClientMessage,
            timestamp: new Date().toISOString(),
            userName: data.userName || (isClientMessage ? ticket.customer.name : 'Agente')
          };
          
          console.log('Mensaje a agregar:', newMessage);
          console.log('===============================================');
          
          // Añadir a la lista de mensajes del ticket
          if (!ticket.messages) {
            ticket.messages = [];
          }
          ticket.messages.push(newMessage);
          
          // Guardar cambios
          saveTickets();
          
          // Broadcastear a todos los clientes conectados
          broadcastMessage(ticketId, newMessage);
          
          res.writeHead(201, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify(newMessage));
        } else {
          res.writeHead(404, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: 'Ticket no encontrado' }));
        }
      } else {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
      }
      return;
    }

    // Widget: enviar mensaje (compatibilidad anterior)
    if (req.method === 'POST' && req.url === '/api/widget/chat/messages') {
      const data = await readBody(req);
      
      // Crear un nuevo ticket
      const newTicket = {
        id: 'TICKET-' + new Date().toISOString().replace(/[^0-9]/g, '').slice(0, 14),
        title: 'Soporte Web - ' + (data.name || 'Anónimo'),
        status: 'open',
        createdAt: new Date().toISOString(),
        customer: {
          name: data.name || 'Anónimo',
          email: data.email || 'anonymous@example.com'
        },
        messages: [
          {
            id: 'MSG-' + Date.now(),
            content: data.description || data.message || '',
            isClient: true, // Siempre es del cliente (Monitor8)
            timestamp: new Date().toISOString()
          }
        ],
        description: data.description || data.message || '',
        priority: 'medium',
        category: 'soporte',
        createdBy: data.email || 'anonymous@example.com',
        assignedTo: null,
        updatedAt: new Date().toISOString()
      };
      
      // Añadir el ticket a la lista
      tickets.push(newTicket);
      
      // Guardar tickets en archivo
      saveTickets();
      
      res.writeHead(201, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({
        id: newTicket.id,
        message: 'Ticket creado con éxito',
        status: 'open'
      }));
      return;
    }

    // Widget: obtener mensajes de un ticket
    if (req.method === 'GET' && req.url.startsWith('/api/widget/chat/messages/')) {
      const ticketId = req.url.split('/').pop();
      const ticket = tickets.find(t => t.id === ticketId);
      
      if (ticket) {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ messages: ticket.messages }));
      } else {
        // Si no hay ticket, devolvemos un error
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Ticket no encontrado' }));
      }
      return;
    }

    // Función para procesar mensajes recibidos del panel admin (sistema de soporte)
    if (req.method === 'POST' && req.url.startsWith('/api/tickets/') && req.url.includes('/messages')) {
      const token = extractToken(req);
      
      if (validateToken(token)) {
        const ticketId = req.url.split('/api/tickets/')[1].split('/messages')[0];
        const data = await readBody(req);
        const ticket = tickets.find(t => t.id === ticketId);
        
        if (ticket) {
          console.log('========== MENSAJE DEL PANEL ADMIN ==========');
          console.log(`URL: ${req.url}`);
          console.log(`Ticket ID: ${ticketId}`);
          console.log(`Contenido: ${data.content}`);
          console.log(`Headers de autenticación: ${req.headers['authorization'] ? 'Presentes' : 'Ausentes'}`);
          console.log(`Headers completos: ${JSON.stringify(req.headers)}`);
          console.log(`Datos completos: ${JSON.stringify(data)}`);
          
          // CORRECCIÓN CRÍTICA: Los mensajes del panel de admin SIEMPRE tienen isClient=false
          // Sin excepciones ni condiciones
          const newMessage = {
            id: 'MSG-' + Date.now(),
            content: data.content,
            isClient: false, // SIEMPRE false para los mensajes del panel admin
            timestamp: new Date().toISOString(),
            created_at: new Date().toISOString(),
            userName: data.userName || 'Agente de Soporte',
            userEmail: data.userEmail || 'soporte@ejemplo.com'
          };
          
          console.log(`Mensaje creado con isClient=${newMessage.isClient}`);
          
          // Añadir a la lista de mensajes del ticket
          if (!ticket.messages) {
            ticket.messages = [];
          }
          ticket.messages.push(newMessage);
          
          // Guardar cambios
          saveTickets();
          
          // Enviar mensaje por WebSocket a todos los clientes
          console.log(`Enviando mensaje a conexiones WebSocket activas...`);
          broadcastMessage(ticketId, newMessage);
          
          res.writeHead(201, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({
            success: true,
            message: 'Mensaje añadido correctamente',
            data: newMessage
          }));
        } else {
          res.writeHead(404, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: 'Ticket no encontrado' }));
        }
      } else {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
      }
      return;
    }

    // Actualizar un ticket (endpoint para actualizar datos o asignar)
    if (req.method === 'PUT' && req.url.startsWith('/api/tickets/')) {
      const token = extractToken(req);
      
      if (validateToken(token)) {
        // Extraer ID del ticket correctamente (eliminando posibles parámetros query)
        let ticketId = req.url.split('/api/tickets/')[1];
        if (ticketId.includes('?')) {
          ticketId = ticketId.split('?')[0];
        }
        
        console.log(`Recibida solicitud PUT para actualizar ticket: ${ticketId}`);
        
        const data = await readBody(req);
        console.log(`Datos recibidos para actualización:`, data);
        
        const ticketIndex = tickets.findIndex(t => t.id === ticketId);
        
        if (ticketIndex !== -1) {
          console.log(`Ticket encontrado en índice ${ticketIndex}, actualizando con datos:`, data);
          
          // Si se está asignando un usuario, actualizar el estado
          if (data.assignedTo !== undefined) {
            if (data.assignedTo) {
              // Si se asigna a alguien, cambiar estado a "assigned"
              data.status = 'assigned';
              console.log(`Asignando ticket a usuario: ${data.assignedTo}, cambiando estado a 'assigned'`);
            } else if (tickets[ticketIndex].assignedTo && !data.assignedTo) {
              // Si se quita la asignación, volver a "open"
              data.status = 'open';
              console.log(`Quitando asignación, cambiando estado a 'open'`);
            }
          }
          
          // Actualizar los campos con los nuevos datos
          const updatedTicket = {
            ...tickets[ticketIndex],
            ...data,
            updatedAt: new Date().toISOString()
          };
          
          // Verificar que el ticket siga teniendo su estructura básica
          if (!updatedTicket.id) {
            updatedTicket.id = ticketId;
          }
          
          // Actualizar el ticket en la lista
          tickets[ticketIndex] = updatedTicket;
          
          // Guardar cambios en el archivo
          saveTickets();
          
          console.log(`Ticket ${ticketId} actualizado correctamente:`, updatedTicket);
          
          // Responder con el ticket actualizado
          res.writeHead(200, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify(updatedTicket));
        } else {
          console.error(`Ticket ${ticketId} no encontrado`);
          res.writeHead(404, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: `Ticket ${ticketId} no encontrado` }));
        }
      } else {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
      }
      return;
    }

    // Obtener usuarios
    if (req.method === 'GET' && req.url === '/api/users') {
      handleGetUsers(req, res);
      return;
    }

    // Obtener usuario por ID
    if (req.method === 'GET' && /^\/api\/users\/[\w-]+$/.test(req.url)) {
      const userId = req.url.split('/').pop();
      const user = users.find(u => u.id === userId);
      
      if (user) {
        const { password, ...safeUser } = user;
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(safeUser));
      } else {
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Usuario no encontrado' }));
      }
      return;
    }

    // Si no coincide con ninguna ruta conocida
    if (path.match(/^\/api\/users\/\d+$/) && method === 'PUT') {
      // Actualizar un usuario existente
      console.log('Detectada solicitud de actualización de usuario:', path);
      const userId = path.split('/').pop();
      console.log('ID de usuario a actualizar:', userId);
      handleUpdateUser(req, res);
      return;
    }

    // === RUTAS PARA CATEGORÍAS ===
    
    // Listar todas las categorías
    if (req.method === 'GET' && req.url === '/api/categories') {
      // Verificar autenticación
      const token = extractToken(req);
      if (!validateToken(token)) {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
        return;
      }
      
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify(categories));
      return;
    }
    
    // Obtener una categoría específica
    if (req.method === 'GET' && req.url.match(/^\/api\/categories\/\d+$/)) {
      // Verificar autenticación
      const token = extractToken(req);
      if (!validateToken(token)) {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
        return;
      }
      
      const id = req.url.split('/').pop();
      const category = categories.find(c => c.id === id);
      
      if (category) {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(category));
      } else {
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Categoría no encontrada' }));
      }
      return;
    }
    
    // Crear una nueva categoría
    if (req.method === 'POST' && req.url === '/api/categories') {
      // Verificar autenticación y rol
      const token = extractToken(req);
      const decodedToken = validateToken(token);
      if (!decodedToken) {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
        return;
      }
      
      // Verificar que el usuario es administrador
      const currentUser = users.find(u => u.id === decodedToken.userId);
      if (!currentUser || currentUser.role !== 'admin') {
        res.writeHead(403, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Acceso denegado. Se requiere rol de administrador.' }));
        return;
      }
      
      // Leer datos de la categoría
      const categoryData = await readBody(req);
      
      // Validar campos requeridos
      if (!categoryData.name) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'El nombre de la categoría es obligatorio' }));
        return;
      }
      
      // Crear nueva categoría
      const newCategory = {
        id: (categories.length + 1).toString(),
        name: categoryData.name,
        description: categoryData.description || '',
        color: categoryData.color || '#2196F3',
        icon: categoryData.icon || 'category',
        active: categoryData.active !== undefined ? categoryData.active : true,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      };
      
      // Añadir a la lista de categorías
      categories.push(newCategory);
      
      // Guardar categorías en archivo
      saveCategories();
      
      // Responder con la categoría creada
      res.writeHead(201, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify(newCategory));
      return;
    }
    
    // Actualizar una categoría existente
    if (req.method === 'PUT' && req.url.match(/^\/api\/categories\/\d+$/)) {
      // Verificar autenticación y rol
      const token = extractToken(req);
      const decodedToken = validateToken(token);
      if (!decodedToken) {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
        return;
      }
      
      // Verificar que el usuario es administrador
      const currentUser = users.find(u => u.id === decodedToken.userId);
      if (!currentUser || currentUser.role !== 'admin') {
        res.writeHead(403, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Acceso denegado. Se requiere rol de administrador.' }));
        return;
      }
      
      // Obtener ID de la categoría
      const id = req.url.split('/').pop();
      
      // Buscar categoría
      const categoryIndex = categories.findIndex(c => c.id === id);
      if (categoryIndex === -1) {
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Categoría no encontrada' }));
        return;
      }
      
      // Leer datos de actualización
      const updateData = await readBody(req);
      
      // Actualizar categoría
      const updatedCategory = {
        ...categories[categoryIndex],
        name: updateData.name || categories[categoryIndex].name,
        description: updateData.description !== undefined ? updateData.description : categories[categoryIndex].description,
        color: updateData.color || categories[categoryIndex].color,
        icon: updateData.icon || categories[categoryIndex].icon,
        active: updateData.active !== undefined ? updateData.active : categories[categoryIndex].active,
        updatedAt: new Date().toISOString()
      };
      
      // Actualizar en el array
      categories[categoryIndex] = updatedCategory;
      
      // Guardar cambios
      saveCategories();
      
      // Responder con la categoría actualizada
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify(updatedCategory));
      return;
    }
    
    // Eliminar una categoría
    if (req.method === 'DELETE' && req.url.match(/^\/api\/categories\/\d+$/)) {
      // Verificar autenticación y rol
      const token = extractToken(req);
      const decodedToken = validateToken(token);
      if (!decodedToken) {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No autorizado' }));
        return;
      }
      
      // Verificar que el usuario es administrador
      const currentUser = users.find(u => u.id === decodedToken.userId);
      if (!currentUser || currentUser.role !== 'admin') {
        res.writeHead(403, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Acceso denegado. Se requiere rol de administrador.' }));
        return;
      }
      
      // Obtener ID de la categoría
      const id = req.url.split('/').pop();
      
      // Buscar categoría
      const categoryIndex = categories.findIndex(c => c.id === id);
      if (categoryIndex === -1) {
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Categoría no encontrada' }));
        return;
      }
      
      // Verificar si la categoría está en uso en algún ticket
      const categoryInUse = tickets.some(ticket => ticket.categoryId === id);
      if (categoryInUse) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ 
          error: 'No se puede eliminar esta categoría porque está siendo utilizada en tickets existentes',
          solution: 'Puede desactivar la categoría en lugar de eliminarla'
        }));
        return;
      }
      
      // Eliminar categoría
      categories.splice(categoryIndex, 1);
      
      // Guardar cambios
      saveCategories();
      
      // Responder éxito
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ message: 'Categoría eliminada correctamente' }));
      return;
    }

    // Si ninguna ruta coincide
    res.writeHead(404, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ error: 'Ruta no encontrada' }));
  } catch (error) {
    console.error('Error interno del servidor:', error);
    res.writeHead(500, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ error: 'Error interno del servidor' }));
  }
});

// Crear servidor WebSocket
const wss = new WebSocket.Server({ noServer: true });

// Manejar actualización a WebSocket
server.on('upgrade', (request, socket, head) => {
  // Verificar si es una ruta de chat
  if (request.url.startsWith('/api/ws/chat/')) {
    const ticketId = request.url.split('/api/ws/chat/')[1];
    console.log(`Solicitando upgrade a WebSocket para ticket: ${ticketId}`);
    console.log('Headers de la solicitud WebSocket:', JSON.stringify(request.headers));
    
    wss.handleUpgrade(request, socket, head, (ws) => {
      // Inicializar la lista de conexiones para este ticket si no existe
      if (!ticketConnections[ticketId]) {
        ticketConnections[ticketId] = [];
      }
      
      // Identificador único para esta conexión
      const connectionId = Date.now().toString();
      
      // Añadir conexión a la lista para este ticket
      ticketConnections[ticketId].push({
        id: connectionId,
        socket: ws,
        connectedAt: new Date()
      });
      
      console.log(`Nueva conexión WebSocket para ticket: ${ticketId}`);
      console.log(`Total de conexiones para este ticket: ${ticketConnections[ticketId].length}`);
      
      // Enviar mensaje de bienvenida
      ws.send(JSON.stringify({
        type: 'connection_established',
        message: {
          id: `system-${Date.now()}`,
          content: 'Conexión establecida',
          isClient: false,
          timestamp: new Date().toISOString()
        }
      }));
      
      // Enviar todos los mensajes actuales para sincronizar
      const ticket = tickets.find(t => t.id === ticketId);
      if (ticket && ticket.messages && ticket.messages.length > 0) {
        console.log(`Enviando ${ticket.messages.length} mensajes existentes al cliente que acaba de conectarse`);
        ws.send(JSON.stringify({
          type: 'message_history',
          messages: ticket.messages,
          ticketId: ticketId
        }));
      }
      
      // Manejar mensajes del cliente
      ws.on('message', async (msg) => {
        try {
          console.log(`Mensaje WebSocket recibido para ticket ${ticketId}:`, msg.toString());
          const data = JSON.parse(msg.toString());
          
          // Verificar tipo de mensaje
          if (data.type === 'message' && data.content) {
            // Crear nuevo mensaje
            const message = {
              id: `MSG-${Date.now()}`,
              ticketId,
              content: data.content,
              // El valor de isClient depende de dónde viene el mensaje:
              // - Si viene de la administración (data.source === 'admin'), debe ser false
              // - Si viene del widget (data.source === 'widget' o no hay source), debe ser true
              isClient: data.source === 'admin' ? false : true,
              userId: data.userId || 'anonymous',
              createdAt: new Date().toISOString(),
              userName: data.userName || (data.source === 'admin' ? 'Agente de Soporte' : 'Cliente'),
              userEmail: data.userEmail || 'user@example.com'
            };
            
            console.log('Mensaje WebSocket procesado:', JSON.stringify(message));
            console.log(`Fuente del mensaje: ${data.source || 'widget'}, isClient: ${message.isClient}`);
            
            // Buscar el ticket
            const ticket = tickets.find(t => t.id === ticketId);
            if (ticket) {
              // Añadir a la lista de mensajes del ticket
              if (!ticket.messages) {
                ticket.messages = [];
              }
              ticket.messages.push(message);
              
              // Actualizar estado del ticket
              ticket.updatedAt = new Date().toISOString();
              
              // Guardar tickets en archivo
              saveTickets();
            }
            
            // Enviar a todos los clientes conectados, incluido el remitente
            broadcastMessage(ticketId, message);
          } else {
            console.error(`Formato de mensaje no reconocido: ${JSON.stringify(data)}`);
            ws.send(JSON.stringify({
              type: 'error',
              message: 'Formato de mensaje no reconocido'
            }));
          }
        } catch (error) {
          console.error(`Error al procesar mensaje de WebSocket: ${error}`);
          ws.send(JSON.stringify({
            type: 'error',
            message: 'Error al procesar el mensaje'
          }));
        }
      });
      
      // Manejar desconexión
      ws.on('close', () => {
        console.log(`Conexión WebSocket cerrada para ticket: ${ticketId}`);
        
        // Eliminar de la lista de conexiones
        if (ticketConnections[ticketId]) {
          ticketConnections[ticketId] = ticketConnections[ticketId].filter(
            conn => conn.id !== connectionId
          );
          
          console.log(`Quedan ${ticketConnections[ticketId].length} conexiones para el ticket ${ticketId}`);
          
          // Si no quedan conexiones, eliminar la entrada
          if (ticketConnections[ticketId].length === 0) {
            delete ticketConnections[ticketId];
            console.log(`Eliminada entrada de conexiones para ticket ${ticketId}`);
          }
        }
      });
      
      // Mantener viva la conexión con ping
      const pingInterval = setInterval(() => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.ping();
        } else {
          clearInterval(pingInterval);
        }
      }, 30000);
    });
  } else {
    socket.destroy();
  }
});

// Función para difundir mensajes a todos los clientes WebSocket
function broadcastMessage(ticketId, message) {
  console.log('======== BROADCAST MENSAJE ========');
  console.log(`Ticket ID: ${ticketId}`);
  console.log(`Mensaje: ${JSON.stringify(message)}`);
  console.log(`isClient: ${message.isClient}`); // Registrar valor de isClient explícitamente
  
  // Buscar conexiones directas para este ticket
  console.log(`Buscando conexiones para ticket ${ticketId}...`);
  
  let connectionFound = false;
  
  // Verificar si hay clientes para este ticket
  if (ticketConnections[ticketId] && ticketConnections[ticketId].length > 0) {
    console.log(`✓ Encontradas ${ticketConnections[ticketId].length} conexiones directas para ticket ${ticketId}`);
    sendToConnections(ticketConnections[ticketId], ticketId, message);
    connectionFound = true;
  } else {
    console.log(`✗ No hay conexiones directas para ticket ${ticketId}`);
  }
  
  // Verificar conexiones alternativas (si está configurada la conexión)
  if (alternateConnectionMap[ticketId]) {
    const altTicketId = alternateConnectionMap[ticketId];
    
    if (ticketConnections[altTicketId] && ticketConnections[altTicketId].length > 0) {
      console.log(`✓ Encontradas ${ticketConnections[altTicketId].length} conexiones alternativas (${altTicketId}) para ticket ${ticketId}`);
      sendToConnections(ticketConnections[altTicketId], ticketId, message);
      connectionFound = true;
    } else {
      console.log(`✗ No hay conexiones alternativas para ticket ${ticketId}`);
    }
  }
  
  if (!connectionFound) {
    console.log(`No se encontraron conexiones para el ticket ${ticketId}. Mensaje en espera.`);
  }
  
  console.log('======== FIN BROADCAST ========');
}

// Función auxiliar para enviar a un conjunto de conexiones
function sendToConnections(connections, ticketId, message) {
  // Preparar mensaje websocket
  // IMPORTANTE: Mantener el isClient original del mensaje
  const wsMessage = {
    type: 'new_message', // Cambiar a minúsculas para coincidencia
    ticketId: ticketId,
    message: {
      id: message.id,
      content: message.content,
      isClient: message.isClient, // MANTENER el valor original de isClient
      timestamp: message.created_at || message.createdAt || message.timestamp || new Date().toISOString(),
      userName: message.userName || (message.isClient ? 'Cliente' : 'Agente de Soporte'),
      userEmail: message.userEmail || 'user@example.com'
    }
  };
  
  // Serializar para envío
  const messageString = JSON.stringify(wsMessage);
  console.log(`Enviando mensaje WebSocket: ${messageString}`);
  console.log(`IMPORTANTE - Valor final de isClient: ${wsMessage.message.isClient}`);
  
  // Enviar a todas las conexiones
  let enviadoCount = 0;
  connections.forEach(client => {
    // Verificar que client sea un objeto WebSocket válido
    if (client && client.socket && client.socket.readyState === WebSocket.OPEN) {
      client.socket.send(messageString);
      enviadoCount++;
    } else if (client && client.readyState === WebSocket.OPEN) {
      client.send(messageString);
      enviadoCount++;
    }
  });
  
  console.log(`Mensaje enviado a ${enviadoCount}/${connections.length} conexiones activas`);
}

// Función para manejar peticiones HTTP
const handleRequest = (req, res) => {
  const { method, url } = req;
  
  // Habilitar CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');
  
  if (method === 'OPTIONS') {
    res.writeHead(200);
    res.end();
    return;
  }
  
  // Extraer ruta base y parámetros
  const urlParts = url.split('?');
  const path = urlParts[0];
  
  // Mostrar petición en consola
  console.log(`${method} ${path}`);
  
  // Ruta para autenticación
  if (path === '/api/auth/login' && method === 'POST') {
    handleLogin(req, res);
    return;
  }
  
  // Rutas para tickets
  if (path === '/api/tickets' && method === 'GET') {
    // Retornar todos los tickets
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify(tickets));
    return;
  }
  
  if (path === '/api/tickets' && method === 'POST') {
    // Crear un nuevo ticket
    handleCreateTicket(req, res);
    return;
  }
  
  if (path.match(/^\/api\/tickets\/\d+$/) && method === 'GET') {
    // Retornar un ticket específico por ID
    const id = path.split('/').pop();
    const ticket = tickets.find(t => t.id === id);
    
    if (ticket) {
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify(ticket));
    } else {
      res.writeHead(404, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Ticket no encontrado' }));
    }
    return;
  }
  
  // Rutas para usuarios
  if (path === '/api/users' && method === 'GET') {
    // Retornar todos los usuarios (sin exponer contraseñas)
    const safeUsers = users.map(user => {
      const { password, ...safeUser } = user;
      return safeUser;
    });
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify(safeUsers));
    return;
  }
  
  if (path === '/api/users' && method === 'POST') {
    // Crear un nuevo usuario
    handleCreateUser(req, res);
    return;
  }
  
  if (path.match(/^\/api\/users\/\d+$/) && method === 'GET') {
    // Retornar un usuario específico por ID
    const id = path.split('/').pop();
    const user = users.find(u => u.id === id);
    
    if (user) {
      const { password, ...safeUser } = user;
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify(safeUser));
    } else {
      res.writeHead(404, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Usuario no encontrado' }));
    }
    return;
  }
  
  if (path.match(/^\/api\/users\/\d+$/) && method === 'PUT') {
    // Actualizar un usuario existente
    console.log('Detectada solicitud de actualización de usuario:', path);
    const userId = path.split('/').pop();
    console.log('ID de usuario a actualizar:', userId);
    handleUpdateUser(req, res);
    return;
  }
  
  if (path.match(/^\/api\/users\/\d+\/role$/) && method === 'PUT') {
    // Actualizar el rol de un usuario
    handleChangeUserRole(req, res);
    return;
  }
  
  if (path.match(/^\/api\/users\/\d+$/) && method === 'DELETE') {
    // Eliminar un usuario
    handleDeleteUser(req, res);
    return;
  }
  
  // Si no coincide con ninguna ruta conocida
  res.writeHead(404, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify({ error: 'Ruta no encontrada' }));
};

// Función para manejar login
const handleLogin = (req, res) => {
  let body = '';
  
  req.on('data', chunk => {
    body += chunk.toString();
  });
  
  req.on('end', () => {
    try {
      const { email, password } = JSON.parse(body);
      const user = users.find(u => u.email === email && u.password === password);
      
      if (user) {
        // En un sistema real, generaríamos un JWT aquí
        const { password, ...userWithoutPassword } = user;
        const token = `mock-token-${user.id}`;
        
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({
          token,
          user: userWithoutPassword
        }));
      } else {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Credenciales inválidas' }));
      }
    } catch (err) {
      console.error('Error en login:', err);
      res.writeHead(400, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Datos de solicitud inválidos' }));
    }
  });
};

// Función para manejar creación de usuario
const handleCreateUser = (req, res) => {
  let body = '';
  
  req.on('data', chunk => {
    body += chunk.toString();
  });
  
  req.on('end', () => {
    try {
      const userData = JSON.parse(body);
      
      // Verificar si el email ya existe
      if (users.some(u => u.email === userData.email)) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'El email ya está registrado' }));
        return;
      }
      
      // Generar ID para el nuevo usuario
      const newId = (users.length > 0) 
        ? String(Math.max(...users.map(u => parseInt(u.id))) + 1) 
        : '1';
      
      // Crear nuevo usuario
      const newUser = {
        id: newId,
        firstName: userData.firstName,
        lastName: userData.lastName,
        email: userData.email,
        password: userData.password, // En un sistema real, esto se hashearía
        role: userData.role || 'employee',
        department: userData.department || null,
        active: userData.active !== undefined ? userData.active : true
      };
      
      users.push(newUser);
      saveUsers();
      
      // Retornar el usuario sin la contraseña
      const { password, ...safeUser } = newUser;
      
      res.writeHead(201, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify(safeUser));
    } catch (err) {
      console.error('Error al crear usuario:', err);
      res.writeHead(400, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Datos de solicitud inválidos' }));
    }
  });
};

// Función para manejar actualización de usuario
const handleUpdateUser = (req, res) => {
  let body = '';
  const id = req.url.split('/').pop();
  
  req.on('data', chunk => {
    body += chunk.toString();
  });
  
  req.on('end', () => {
    try {
      const userData = JSON.parse(body);
      const userIndex = users.findIndex(u => u.id === id);
      
      if (userIndex === -1) {
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Usuario no encontrado' }));
        return;
      }
      
      // Verificar si está intentando actualizar el email a uno que ya existe
      if (userData.email && userData.email !== users[userIndex].email && 
          users.some(u => u.email === userData.email)) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'El email ya está registrado por otro usuario' }));
        return;
      }
      
      // Verificar permisos - solo administradores y asistentes pueden modificar otros usuarios
      const authHeader = req.headers.authorization;
      const token = authHeader && authHeader.split(' ')[1];
      const decodedToken = validateToken(token);
      
      if (!decodedToken) {
        res.writeHead(401, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Token inválido o expirado' }));
        return;
      }
      
      const currentUser = users.find(u => u.id === decodedToken.userId);
      
      // Si no es admin ni asistente y está intentando modificar a otro usuario
      if (currentUser.role !== 'admin' && currentUser.role !== 'assistant' && id !== currentUser.id) {
        res.writeHead(403, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'No tienes permisos para modificar este usuario' }));
        return;
      }
      
      // Solo los administradores pueden cambiar roles
      if (userData.role && userData.role !== users[userIndex].role && currentUser.role !== 'admin') {
        res.writeHead(403, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Solo los administradores pueden cambiar roles' }));
        return;
      }
      
      // Actualizar campos del usuario
      const updatedUser = {
        ...users[userIndex],
        firstName: userData.firstName || users[userIndex].firstName,
        lastName: userData.lastName || users[userIndex].lastName,
        email: userData.email || users[userIndex].email,
        department: userData.department !== undefined ? userData.department : users[userIndex].department,
        active: userData.active !== undefined ? userData.active : users[userIndex].active,
        // Campos adicionales del perfil
        position: userData.position !== undefined ? userData.position : users[userIndex].position,
        phone: userData.phone !== undefined ? userData.phone : users[userIndex].phone,
        language: userData.language || users[userIndex].language || 'es',
      };
      
      // Actualizar rol si viene en la solicitud y el usuario tiene permisos
      if (userData.role && currentUser.role === 'admin') {
        updatedUser.role = userData.role;
      }
      
      // Actualizar contraseña solo si viene en la solicitud
      if (userData.password) {
        updatedUser.password = userData.password;
      }
      
      users[userIndex] = updatedUser;
      saveUsers();
      
      // Retornar el usuario actualizado sin la contraseña
      const { password, ...safeUser } = updatedUser;
      
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify(safeUser));
    } catch (err) {
      console.error('Error al actualizar usuario:', err);
      res.writeHead(400, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Datos de solicitud inválidos' }));
    }
  });
};

// Función para manejar cambio de rol de usuario
const handleChangeUserRole = (req, res) => {
  let body = '';
  const id = req.url.split('/')[3]; // Extrae el ID de la URL /api/users/{id}/role
  
  req.on('data', chunk => {
    body += chunk.toString();
  });
  
  req.on('end', () => {
    try {
      const { role } = JSON.parse(body);
      const userIndex = users.findIndex(u => u.id === id);
      
      if (userIndex === -1) {
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Usuario no encontrado' }));
        return;
      }
      
      // Validar rol
      if (!['admin', 'assistant', 'employee'].includes(role)) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Rol no válido' }));
        return;
      }
      
      // Actualizar rol del usuario
      users[userIndex].role = role;
      saveUsers();
      
      // Retornar el usuario actualizado sin la contraseña
      const { password, ...safeUser } = users[userIndex];
      
      res.writeHead(200, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify(safeUser));
    } catch (err) {
      console.error('Error al cambiar rol de usuario:', err);
      res.writeHead(400, { 'Content-Type': 'application/json' });
      res.end(JSON.stringify({ error: 'Datos de solicitud inválidos' }));
    }
  });
};

// Función para manejar eliminación de usuario
const handleDeleteUser = (req, res) => {
  const id = req.url.split('/').pop();
  const userIndex = users.findIndex(u => u.id === id);
  
  if (userIndex === -1) {
    res.writeHead(404, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ error: 'Usuario no encontrado' }));
    return;
  }
  
  // Eliminar usuario
  users.splice(userIndex, 1);
  saveUsers();
  
  res.writeHead(200, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify({ message: 'Usuario eliminado correctamente' }));
};

// Controlador para obtener todos los usuarios
const handleGetUsers = (req, res) => {
  // Eliminar la contraseña de los usuarios para la respuesta
  const safeUsers = users.map(user => {
    const { password, ...safeUser } = user;
    return safeUser;
  });

  res.writeHead(200, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify(safeUsers));
};

// Iniciar el servidor
server.listen(PORT, () => {
  console.log(`Servidor mock ejecutándose en http://localhost:${PORT}`);
});