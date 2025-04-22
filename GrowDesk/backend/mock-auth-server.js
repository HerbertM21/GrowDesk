const http = require('http');

const PORT = process.env.PORT || 8000;

// Token simple para pruebas
const generateToken = () => {
  return 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJhZG1pbi0xMjMiLCJlbWFpbCI6ImFkbWluQGdyb3dkZXNrLmNvbSIsInJvbGUiOiJhZG1pbiIsImV4cCI6MTcyNDA4ODQwMH0.8J5ayPvA4B-1vF5NaqRXycW1pmIl9qjKP6Ddj4Ot_Cw';
};

// Una función para extraer el token de authorization
const extractToken = (req) => {
  const authHeader = req.headers['authorization'];
  if (authHeader && authHeader.startsWith('Bearer ')) {
    return authHeader.substring(7);
  }
  return null;
};

// Validar el token (en este ejemplo simplificado, solo verificamos que exista)
const validateToken = (token) => {
  return token && token.length > 0;
};

// Base de datos en memoria para tickets
const tickets = [
  {
    id: 'TICKET-20250321-230744',
    title: 'Test message',
    status: 'open',
    createdAt: '2025-03-21T23:07:44-03:00',
    customer: {
      name: 'Test User',
      email: 'test@example.com'
    },
    messages: [
      {
        id: 'MSG-1',
        content: 'Test message',
        isClient: true,
        timestamp: '2025-03-21T23:07:44-03:00'
      }
    ]
  }
];

const server = http.createServer((req, res) => {
  // Configurar cabeceras CORS para permitir solicitudes desde cualquier origen
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-Requested-With, X-Widget-ID, X-Widget-Token');
  res.setHeader('Access-Control-Allow-Credentials', 'true');
  res.setHeader('Access-Control-Max-Age', '86400'); // 24 horas
  
  // Manejar solicitudes OPTIONS preflight
  if (req.method === 'OPTIONS') {
    res.writeHead(204);
    res.end();
    return;
  }

  // Servir rutas API
  if (req.method === 'GET' && req.url === '/api/health') {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    res.end(JSON.stringify({ status: 'ok', message: 'Server running' }));
    return;
  }

  // Endpoint para obtener información del usuario autenticado
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

  // Endpoint para obtener tickets
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

  // Manejar login
  if (req.method === 'POST' && req.url === '/api/auth/login') {
    let body = '';
    req.on('data', chunk => {
      body += chunk.toString();
    });
    
    req.on('end', () => {
      try {
        const data = JSON.parse(body);
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
      } catch (e) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Datos inválidos' }));
      }
    });
    return;
  }

  // Manejar registro
  if (req.method === 'POST' && req.url === '/api/auth/register') {
    let body = '';
    req.on('data', chunk => {
      body += chunk.toString();
    });
    
    req.on('end', () => {
      try {
        const data = JSON.parse(body);
        if (data.email && data.password && data.firstName && data.lastName) {
          const token = generateToken();
          
          res.writeHead(201, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({
            token: token,
            user: {
              id: 'user-' + Date.now(),
              email: data.email,
              firstName: data.firstName,
              lastName: data.lastName,
              role: 'customer'
            }
          }));
        } else {
          res.writeHead(400, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: 'Todos los campos son requeridos' }));
        }
      } catch (e) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Datos inválidos' }));
      }
    });
    return;
  }
  
  // Manejar widget routes
  if (req.url.startsWith('/api/widget/')) {
    res.writeHead(200, { 'Content-Type': 'application/json' });
    
    if (req.url === '/api/widget/tickets' && req.method === 'POST') {
      // Capturar los datos del ticket
      let body = '';
      req.on('data', chunk => {
        body += chunk.toString();
      });
      
      req.on('end', () => {
        try {
          const data = JSON.parse(body);
          const newTicket = {
            id: 'TICKET-' + new Date().toISOString().replace(/[^0-9]/g, '').slice(0, 14),
            title: data.message || 'Sin asunto',
            status: 'open',
            createdAt: new Date().toISOString(),
            customer: {
              name: data.name || 'Anónimo',
              email: data.email || 'anonymous@example.com'
            },
            messages: [
              {
                id: 'MSG-' + Date.now(),
                content: data.message || '',
                isClient: true,
                timestamp: new Date().toISOString()
              }
            ]
          };
          
          // Añadir el ticket a la lista
          tickets.push(newTicket);
          
          // Responder
          res.writeHead(201, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({
            success: true,
            ticketId: newTicket.id,
            liveChatAvailable: true
          }));
        } catch (e) {
          res.writeHead(400, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: 'Datos inválidos' }));
        }
      });
      return;
    }
    
    if (req.url === '/api/widget/messages' && req.method === 'POST') {
      let body = '';
      req.on('data', chunk => {
        body += chunk.toString();
      });
      
      req.on('end', () => {
        try {
          const data = JSON.parse(body);
          
          res.writeHead(201, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({
            success: true,
            messageId: 'MSG-' + Date.now(),
            timestamp: new Date().toISOString()
          }));
        } catch (e) {
          res.writeHead(400, { 'Content-Type': 'application/json' });
          res.end(JSON.stringify({ error: 'Datos inválidos' }));
        }
      });
      return;
    }
    
    if (req.url.startsWith('/api/widget/tickets/') && req.method === 'GET') {
      const ticketId = req.url.split('/').pop();
      const ticket = tickets.find(t => t.id === ticketId);
      
      if (ticket) {
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(ticket.messages));
      } else {
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Ticket no encontrado' }));
      }
      return;
    }
  }

  // Ruta no encontrada
  res.writeHead(404, { 'Content-Type': 'application/json' });
  res.end(JSON.stringify({ error: 'Ruta no encontrada' }));
});

server.listen(PORT, () => {
  console.log(`Mock servidor ejecutándose en http://localhost:${PORT}`);
}); 