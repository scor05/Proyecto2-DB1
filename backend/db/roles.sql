/*
    Roles y permisos para Proyecto 3.
*/

CREATE ROLE rol_cliente;
CREATE ROLE rol_empleado;
CREATE ROLE rol_proveedor;
CREATE ROLE rol_gerente;
CREATE ROLE rol_superadmin;

REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public FROM PUBLIC;
REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public FROM PUBLIC;

REVOKE ALL PRIVILEGES ON ALL TABLES IN SCHEMA public
FROM rol_cliente, rol_empleado, rol_proveedor, rol_gerente, rol_superadmin;

REVOKE ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public
FROM rol_cliente, rol_empleado, rol_proveedor, rol_gerente, rol_superadmin;

REVOKE ALL PRIVILEGES ON SCHEMA public
FROM rol_cliente, rol_empleado, rol_proveedor, rol_gerente, rol_superadmin;

GRANT CONNECT ON DATABASE proyecto2
TO rol_cliente, rol_empleado, rol_proveedor, rol_gerente, rol_superadmin;

GRANT USAGE ON SCHEMA public
TO rol_cliente, rol_empleado, rol_proveedor, rol_gerente;

GRANT LOGIN, USAGE, CREATE ON SCHEMA public TO rol_superadmin;

/*
    rol_cliente:
    Usuario externo que puede consultar el catalogo publico.
*/
GRANT SELECT ON TABLE categoria, producto TO rol_cliente;

/*
    rol_empleado:
    Usuario operativo de ventas. Puede consultar catalogo, clientes y compras;
    registrar clientes/compras; y ajustar stock durante operaciones de venta.
*/
GRANT SELECT ON TABLE categoria, proveedor, producto, cliente, empleado, compra, producto_compra
TO rol_empleado;

GRANT INSERT, UPDATE ON TABLE cliente TO rol_empleado;
GRANT INSERT, UPDATE ON TABLE compra TO rol_empleado;
GRANT INSERT, UPDATE, DELETE ON TABLE producto_compra TO rol_empleado;
GRANT UPDATE (stock) ON TABLE producto TO rol_empleado;
GRANT USAGE, SELECT ON SEQUENCE cliente_id_cliente_seq, compra_id_compra_seq TO rol_empleado;

/*
    rol_proveedor:
    Usuario externo de proveedor. Puede consultar catalogo y proveedores;
    solo puede actualizar datos de contacto de proveedor.
*/
GRANT SELECT ON TABLE categoria, producto, proveedor TO rol_proveedor;
GRANT UPDATE (telefono, correo, direccion) ON TABLE proveedor TO rol_proveedor;

/*
    rol_gerente:
    Usuario de administracion del negocio. Puede consultar toda la informacion
    operativa y ajustar precio/stock de productos, pero no administra el sistema.
*/
GRANT SELECT ON TABLE categoria, proveedor, empleado, cliente, producto, compra, producto_compra, gerente
TO rol_gerente;

GRANT UPDATE (precio, stock) ON TABLE producto TO rol_gerente;

/*
    rol_superadmin:
    Usuario para desarrolladores con control completo sobre datos de la app.
*/
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO rol_superadmin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO rol_superadmin;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT ALL PRIVILEGES ON TABLES TO rol_superadmin;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT ALL PRIVILEGES ON SEQUENCES TO rol_superadmin;
