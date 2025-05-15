import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/auth/Login.vue'
import Dashboard from '../views/dashboard/Dashboard.vue'
import { useAuthStore } from '../stores/auth'
import { useNotificationsStore } from '../stores/notifications'
import NotFound from '../views/NotFound.vue'

// Definir las rutas
const routes = [
  {
    path: '/',
    name: 'home',
    component: Login,
    meta: {
      requiresAuth: false
    }
  },
  {
    path: '/login',
    name: 'login',
    component: Login,
    meta: {
      requiresAuth: false
    }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/auth/Register.vue')
  },
  {
    path: '/tickets',
    name: 'tickets',
    component: () => import('../views/tickets/TicketList.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/tickets-board',
    name: 'tickets-board',
    component: () => import('../views/tickets/KanbanBoard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/tickets/:id',
    name: 'ticket-detail',
    component: () => import('../views/tickets/TicketDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: Dashboard,
    meta: {
      requiresAuth: true
    }
  },
  {
    path: '/profile',
    name: 'user-profile',
    component: () => import('../views/Profile/UserProfile.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'user-settings',
    component: () => import('../views/Profile/UserSettings.vue'),
    meta: { requiresAuth: true }
  },
  // Rutas de administración
  {
    path: "/admin",
    redirect: "/admin/dashboard",
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: "/admin/dashboard",
    name: "admin-dashboard",
    component: () => import("../views/Admin/AdminDashboard.vue"),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: "/admin/users",
    name: "admin-users",
    component: () => import("../views/Admin/UsersList.vue"),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: "/admin/profile-management",
    name: "admin-profile-management",
    component: () => import("../views/Admin/ProfileManagement.vue"),
    meta: { requiresAuth: true, requiresAdminOrAssistant: true }
  },
  {
    path: "/admin/categories",
    name: "admin-categories",
    component: () => import("../views/Admin/CategoriesManagement.vue"),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: "/admin/faqs",
    name: "admin-faqs",
    component: () => import("../views/Admin/FaqManagement.vue"),
    meta: { requiresAuth: true, requiresAdminOrAssistant: true }
  },
  {
    path: "/admin/widget-config",
    name: "admin-widget-config",
    component: () => import("../views/Admin/WidgetConfigView.vue"),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  // Ruta para manejar errores 404
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFound,
    meta: { requiresAuth: false }
  }
]

// Crear la instancia del router
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
  // Scroll behavior cuando cambia la ruta
  scrollBehavior() {
    // Siempre scroll al principio de la página
    return { top: 0 }
  }
})

// Protección de rutas
router.beforeEach((to: any, from: any, next: any) => {
  const authStore = useAuthStore()
  const notificationsStore = useNotificationsStore()
  
  // Comprobar si el usuario está autenticado
  const isAuthenticated = authStore.isAuthenticated
  
  // Si la ruta requiere autenticación y el usuario no está autenticado
  if (to.meta.requiresAuth && !isAuthenticated) {
    notificationsStore.error('Debes iniciar sesión para acceder a esta página')
    next({ name: 'login' })
  } 
  // Si la ruta es la página principal y el usuario ya está autenticado
  else if ((to.path === '/' || to.path === '/login') && isAuthenticated) {
    // Redirigir a dashboard o admin según el rol
    if (authStore.isAdmin || authStore.isAssistant) {
      next({ name: 'admin-dashboard' })
    } else {
      next({ name: 'dashboard' })
    }
  }
  // Si la ruta requiere rol de admin y el usuario no es admin
  else if (to.meta.requiresAdmin && !authStore.isAdmin) {
    notificationsStore.error('No tienes permisos suficientes para acceder a esta sección')
    next({ name: 'dashboard' })
  }
  // Si la ruta requiere rol de admin o asistente y el usuario no es ninguno de esos roles
  else if (to.meta.requiresAdminOrAssistant && !(authStore.isAdmin || authStore.isAssistant)) {
    notificationsStore.error('No tienes permisos suficientes para acceder a esta sección')
    next({ name: 'dashboard' })
  }
  // En otros casos, permitir la navegación
  else {
    next()
  }
})

// Exportar la instancia del router
export default router