<!-- Sección tickets -->
      <div class="sidebar-section">
        <h3 class="sidebar-section-title">Tickets</h3>
        <ul>
          <li>
            <router-link to="/tickets" class="sidebar-link" active-class="active">
              <i class="pi pi-ticket"></i>
              <span>Todos los tickets</span>
            </router-link>
          </li>
          <li>
            <router-link to="/my-tickets" class="sidebar-link" active-class="active">
              <i class="pi pi-user-edit"></i>
              <span>Mis tickets asignados</span>
              <span v-if="ticketStore.assignedTickets.length > 0" class="ticket-count">
                {{ ticketStore.assignedTickets.length }}
              </span>
            </router-link>
          </li>
          <li v-if="authStore.isAdmin">
            <router-link to="/tickets/new" class="sidebar-link" active-class="active">
              <i class="pi pi-plus-circle"></i>
              <span>Crear ticket</span>
            </router-link>
          </li>
        </ul>
      </div> 