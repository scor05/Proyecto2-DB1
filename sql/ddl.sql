/*
    Archivo DDL para el sistema de Proyecto
    Autor: Santiago Cordero
    Carnet: 24472
*/

CREATE TABLE categoria (
    id_categoria INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    nombre VARCHAR(100) NOT NULL
);

CREATE TABLE proveedor (
    id_proveedor INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    nombre VARCHAR(150) NOT NULL,
    telefono VARCHAR(20) NOT NULL,
    correo VARCHAR(250) NOT NULL,
    direccion VARCHAR(250) NOT NULL
);

CREATE TABLE empleado (
    id_empleado INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    nombre VARCHAR(150) NOT NULL,
    estado VARCHAR(8) NOT NULL CHECK (estado IN ('activo', 'inactivo')),
    correo VARCHAR(250) NOT NULL
);

CREATE TABLE cliente (
    id_cliente INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    nombre VARCHAR(150) NOT NULL,
    telefono VARCHAR(20) NOT NULL,
    correo VARCHAR(254) NOT NULL
);

CREATE TABLE producto (
    id_producto INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    id_categoria INTEGER NOT NULL,
    id_proveedor INTEGER NOT NULL,
    precio DECIMAL(10, 2) NOT NULL CHECK (precio > 0),
    stock INTEGER NOT NULL CHECK (stock >= 0),
    nombre VARCHAR(200) NOT NULL,
    imagen TEXT,
    descripcion TEXT NOT NULL,
    FOREIGN KEY (id_categoria) REFERENCES categoria (id_categoria),
    FOREIGN KEY (id_proveedor) REFERENCES proveedor (id_proveedor)
);

CREATE TABLE compra (
    id_compra INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    id_empleado INTEGER NOT NULL,
    id_cliente INTEGER NOT NULL,
    fecha_compra DATE NOT NULL,
    FOREIGN KEY (id_empleado) REFERENCES empleado (id_empleado),
    FOREIGN KEY (id_cliente) REFERENCES cliente (id_cliente)
);

CREATE TABLE producto_compra (
    id_compra INTEGER NOT NULL,
    id_producto INTEGER NOT NULL,
    cantidad_producto INTEGER NOT NULL CHECK (cantidad_producto > 0),
    PRIMARY KEY (id_producto, id_compra),
    FOREIGN KEY (id_compra) REFERENCES compra (id_compra),
    FOREIGN KEY (id_producto) REFERENCES producto (id_producto)
);

/*
    Casi que es obligatorio poner un índice en categoría para un sistema como este, pues las consultas del catalogo para filtrar por productos (como buscar solo GPUs, CPUs o demás) es algo muy frecuente y casi la funcionalidad principal del proyecto.
*/ 
CREATE INDEX idx_producto_id_categoria ON producto (id_categoria);

/*
    Principalmente este índice está para los gerentes y empleados de la tienda, para hacer más eficiente la generación de reportes y busquedas/filtros de compras por fecha, lo cual es algo que también se hace bastante frecentemente en casi cualquier sistema.
*/
CREATE INDEX idx_compra_fecha_compra ON compra (fecha_compra);
