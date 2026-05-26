/*
    Procedimientos para productos.
*/

CREATE OR REPLACE PROCEDURE sp_create_producto(
    IN p_id_producto INTEGER,
    IN p_id_categoria INTEGER,
    IN p_id_proveedor INTEGER,
    IN p_precio NUMERIC,
    IN p_stock INTEGER,
    IN p_nombre TEXT,
    IN p_imagen TEXT,
    IN p_descripcion TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO producto (
        id_producto,
        id_categoria,
        id_proveedor,
        precio,
        stock,
        nombre,
        imagen,
        descripcion
    )
    OVERRIDING SYSTEM VALUE
    VALUES (
        p_id_producto,
        p_id_categoria,
        p_id_proveedor,
        p_precio,
        p_stock,
        p_nombre,
        p_imagen,
        p_descripcion
    );
END;
$$;

CREATE OR REPLACE PROCEDURE sp_update_producto(
    IN p_id_producto INTEGER,
    IN p_id_categoria INTEGER,
    IN p_id_proveedor INTEGER,
    IN p_precio NUMERIC,
    IN p_stock INTEGER,
    IN p_nombre TEXT,
    IN p_imagen TEXT,
    IN p_descripcion TEXT
)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE producto
    SET id_categoria = p_id_categoria,
        id_proveedor = p_id_proveedor,
        precio = p_precio,
        stock = p_stock,
        nombre = p_nombre,
        imagen = p_imagen,
        descripcion = p_descripcion
    WHERE id_producto = p_id_producto;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'producto_not_found';
    END IF;
END;
$$;

CREATE OR REPLACE PROCEDURE sp_delete_producto(
    IN p_id_producto INTEGER
)
LANGUAGE plpgsql
AS $$
BEGIN
    DELETE FROM producto_compra
    WHERE id_producto = p_id_producto;

    DELETE FROM producto
    WHERE id_producto = p_id_producto;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'producto_not_found';
    END IF;
END;
$$;
