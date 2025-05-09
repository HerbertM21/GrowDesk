// Script para forzar el acceso de administrador
(function() {
  // Crear el token mock
  const token = `mock-jwt-token-${Date.now()}`;
  
  // Crear ID de usuario administrador
  const userId = "1";
  
  // Crear usuario administrador
  const mockUsers = [
    {
      id: "1",
      email: "admin@example.com",
      firstName: "Herbert",
      lastName: "Usuario",
      role: "admin",
      department: "Tecnología",
      active: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      position: "Gerente de TI",
      phone: "+569 1234 5678",
      language: "es"
    },
    {
      id: "2",
      email: "asistente@example.com",
      firstName: "Asistente",
      lastName: "Soporte",
      role: "assistant",
      department: "Soporte",
      active: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      position: "Coordinador de Soporte",
      phone: "+34 600 234 567",
      language: "es"
    },
    {
      id: "3",
      email: "empleado@example.com",
      firstName: "Empleado",
      lastName: "Regular",
      role: "employee",
      department: "Ventas",
      active: true,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      position: "Representante de Ventas",
      phone: "+34 600 345 678",
      language: "en"
    }
  ];
  
  // Guardar datos en localStorage
  localStorage.setItem('token', token);
  localStorage.setItem('userId', userId);
  localStorage.setItem('growdesk-users', JSON.stringify(mockUsers));
  
  console.log('Datos de administrador inyectados con éxito!');
  console.log('Token:', token);
  console.log('UserId:', userId);
  
  // Redireccionar al dashboard de administrador
  if (window.location.pathname === '/login') {
    window.location.href = '/admin/dashboard';
  }
})(); 