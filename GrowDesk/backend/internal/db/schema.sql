-- schema.sql - Definición del esquema de la base de datos GrowDesk

-- Extensiones necesarias
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Tabla de usuarios
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role TEXT NOT NULL CHECK (role IN ('admin', 'assistant', 'employee', 'customer')),
    department TEXT,
    active BOOLEAN DEFAULT TRUE,
    position TEXT,
    phone TEXT,
    language TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de categorías
CREATE TABLE IF NOT EXISTS categories (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    color TEXT,
    icon TEXT,
    active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de tickets
CREATE TABLE IF NOT EXISTS tickets (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    subject TEXT,
    description TEXT,
    status TEXT NOT NULL CHECK (status IN ('open', 'pending', 'in_progress', 'resolved', 'closed')),
    priority TEXT CHECK (priority IN ('low', 'medium', 'high', 'urgent')),
    category TEXT,
    category_id TEXT REFERENCES categories(id),
    assigned_to TEXT REFERENCES users(id),
    created_by TEXT REFERENCES users(id),
    user_id TEXT REFERENCES users(id),
    source TEXT,
    widget_id TEXT,
    department TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de metadatos de tickets
CREATE TABLE IF NOT EXISTS ticket_metadata (
    id SERIAL PRIMARY KEY,
    ticket_id TEXT REFERENCES tickets(id) ON DELETE CASCADE,
    url TEXT,
    referrer TEXT,
    user_agent TEXT,
    screen_size TEXT,
    external_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de clientes de tickets
CREATE TABLE IF NOT EXISTS customers (
    id SERIAL PRIMARY KEY,
    ticket_id TEXT REFERENCES tickets(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de mensajes
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    ticket_id TEXT REFERENCES tickets(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_client BOOLEAN DEFAULT FALSE,
    is_internal BOOLEAN DEFAULT FALSE,
    user_id TEXT REFERENCES users(id),
    user_name TEXT,
    user_email TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de FAQs
CREATE TABLE IF NOT EXISTS faqs (
    id SERIAL PRIMARY KEY,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    category TEXT,
    is_published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de configuración de widgets
CREATE TABLE IF NOT EXISTS widget_settings (
    id TEXT PRIMARY KEY,
    widget_id TEXT UNIQUE NOT NULL,
    widget_token TEXT UNIQUE NOT NULL,
    brand_name TEXT NOT NULL,
    welcome_message TEXT,
    primary_color TEXT,
    position TEXT CHECK (position IN ('left', 'right')),
    logo_url TEXT,
    allowed_domains JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de tickets de widget
CREATE TABLE IF NOT EXISTS widget_tickets (
    id TEXT PRIMARY KEY,
    ticket_id TEXT REFERENCES tickets(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    subject TEXT,
    description TEXT,
    status TEXT NOT NULL,
    priority TEXT,
    client_name TEXT NOT NULL,
    client_email TEXT NOT NULL,
    widget_id TEXT REFERENCES widget_settings(widget_id),
    department TEXT,
    synced BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de mensajes de widget
CREATE TABLE IF NOT EXISTS widget_messages (
    id TEXT PRIMARY KEY,
    widget_ticket_id TEXT REFERENCES widget_tickets(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_client BOOLEAN DEFAULT FALSE,
    user_name TEXT,
    user_email TEXT,
    synced BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de sesiones de widget
CREATE TABLE IF NOT EXISTS widget_sessions (
    id TEXT PRIMARY KEY,
    name TEXT,
    email TEXT,
    ticket_id TEXT REFERENCES widget_tickets(id) ON DELETE CASCADE,
    widget_id TEXT REFERENCES widget_settings(widget_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE,
    last_active TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de metadatos de widget
CREATE TABLE IF NOT EXISTS widget_metadata (
    id TEXT PRIMARY KEY,
    session_id TEXT REFERENCES widget_sessions(id) ON DELETE CASCADE,
    url TEXT,
    referrer TEXT,
    user_agent TEXT,
    screen_size TEXT,
    browser TEXT,
    os TEXT,
    ip_address TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de actividades
CREATE TABLE IF NOT EXISTS activities (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id),
    type TEXT NOT NULL,
    target_id TEXT,
    description TEXT NOT NULL,
    metadata JSONB,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de notificaciones
CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    type TEXT NOT NULL,
    read BOOLEAN DEFAULT FALSE,
    related_id TEXT,
    related_type TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Tabla de adjuntos
CREATE TABLE IF NOT EXISTS attachments (
    id TEXT PRIMARY KEY,
    message_id TEXT REFERENCES messages(id) ON DELETE CASCADE,
    file_name TEXT NOT NULL,
    file_type TEXT NOT NULL,
    file_size INTEGER NOT NULL,
    file_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Índices
CREATE INDEX IF NOT EXISTS idx_tickets_status ON tickets(status);
CREATE INDEX IF NOT EXISTS idx_tickets_user_id ON tickets(user_id);
CREATE INDEX IF NOT EXISTS idx_tickets_category_id ON tickets(category_id);
CREATE INDEX IF NOT EXISTS idx_tickets_assigned_to ON tickets(assigned_to);
CREATE INDEX IF NOT EXISTS idx_messages_ticket_id ON messages(ticket_id);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_faqs_is_published ON faqs(is_published);
CREATE INDEX IF NOT EXISTS idx_widget_tickets_widget_id ON widget_tickets(widget_id);
CREATE INDEX IF NOT EXISTS idx_widget_messages_widget_ticket_id ON widget_messages(widget_ticket_id);
CREATE INDEX IF NOT EXISTS idx_activities_user_id ON activities(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(read); 