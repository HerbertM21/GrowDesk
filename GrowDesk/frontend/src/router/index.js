// Importar componentes
import Home from '../views/Home.vue'
import LoginView from '../views/auth/Login.vue'
import RegisterView from '../views/auth/Register.vue'
import TicketList from '../views/tickets/TicketList.vue'
import TicketDetail from '../views/tickets/TicketDetail.vue'
import CreateTicket from '../views/tickets/CreateTicket.vue'
import AssignedTickets from '../views/tickets/AssignedTickets.vue'
import ProfileView from '../views/Profile/ProfileView.vue'
import { useAuthStore } from '../stores/auth'

// Otras importaciones...

// Rutas
const routes = [
  // Ruta principal y autenticación
  {
    path: '/',
    name: 'home',
    component: Home,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: { guestOnly: true }
  },
  {
    path: '/register',
    name: 'register',
    component: RegisterView,
    meta: { guestOnly: true }
  },
  
  // Rutas de tickets
  {
    path: '/tickets',
    name: 'tickets',
    component: TicketList,
    meta: { requiresAuth: true }
  },
  {
    path: '/my-tickets',
    name: 'assignedTickets',
    component: AssignedTickets,
    meta: { requiresAuth: true }
  },
  {
    path: '/tickets/new',
    name: 'createTicket',
    component: CreateTicket,
    meta: { requiresAuth: true }
  },
  {
    path: '/tickets/:id',
    name: 'ticketDetail',
    component: TicketDetail,
    meta: { requiresAuth: true }
  },
  
  // Perfil de usuario
  {
    path: '/profile',
    name: 'profile',
    component: ProfileView,
    meta: { requiresAuth: true }
  },
  
  // Otras rutas...
] 