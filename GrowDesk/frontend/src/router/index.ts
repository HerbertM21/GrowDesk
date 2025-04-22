import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/auth/Login.vue'
import Dashboard from '../views/dashboard/Dashboard.vue'
import { useAuthStore } from '../stores/auth'

// Definir las rutas
const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/Home.vue'),
    meta: { requiresAuth: false }
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
    path: '/tickets/new',
    name: 'new-ticket',
    component: () => import('../views/tickets/CreateTicket.vue'),
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
    path: '/admin',
    name: 'admin',
    component: () => import('../views/Admin/UsersList.vue'),
    meta: {
      requiresAuth: true,
      requiresAdmin: true
    }
  },
  {
    path: '/admin/widget-config',
    name: 'widget-config',
    component: () => import('../views/Admin/WidgetConfigView.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/admin/users',
    name: 'admin-users',
    component: () => import('../views/Admin/UsersList.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/admin/profile-management',
    name: 'profile-management',
    component: () => import('../views/Admin/ProfileManagement.vue'),
    meta: { requiresAuth: true, requiresAdminOrAssistant: true }
  },
  {
    path: '/admin/categories',
    name: 'categories-management',
    component: () => import('../views/Admin/CategoriesManagement.vue'),
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: '/admin/faqs',
    name: 'faqs-management',
    component: () => import('../views/Admin/FaqManagement.vue'),
    meta: { requiresAuth: true, requiresAdminOrAssistant: true }
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
  }
]

// Rutas de administración
const adminRoutes = [
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
    meta: { requiresAuth: true, requiresAdmin: true }
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
    meta: { requiresAuth: true, requiresAdmin: true }
  },
  {
    path: "/admin/widget-config",
    name: "admin-widget-config",
    component: () => import("../views/Admin/WidgetConfigView.vue"),
    meta: { requiresAuth: true, requiresAdmin: true }
  }
];

// Añadir rutas de administración
routes.push(...adminRoutes);

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
// @ts-expect-error 
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // Si la ruta requiere autenticación y el usuario no está autenticado
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login' })
  } 
  // Si la ruta requiere rol de admin y el usuario no es admin
  else if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'dashboard' })
  }
  // En otros casos, permitir la navegación
  else {
    next()
  }
})

// Exportar la instancia del router
export default router 