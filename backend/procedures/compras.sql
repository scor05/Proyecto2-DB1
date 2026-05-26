/*
    Procedimientos para compras.
*/

CREATE OR REPLACE PROCEDURE sp_create_compra(
    IN p_id_compra INTEGER,
    IN p_id_empleado INTEGER,
    IN p_id_cliente INTEGER,
    IN p_fecha_compra DATE,
    IN p_total_compra NUMERIC,
    IN p_id_producto INTEGER,
    IN p_cantidad_producto INTEGER
)
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO compra (
        id_compra,
        id_empleado,
        id_cliente,
        fecha_compra,
        total_compra
    )
    OVERRIDING SYSTEM VALUE
    VALUES (
        p_id_compra,
        p_id_empleado,
        p_id_cliente,
        p_fecha_compra,
        p_total_compra
    );

    INSERT INTO producto_compra (
        id_compra,
        id_producto,
        cantidad_producto
    )
    VALUES (
        p_id_compra,
        p_id_producto,
        p_cantidad_producto
    );
END;
$$;

CREATE OR REPLACE PROCEDURE sp_update_compra(
    IN p_id_compra INTEGER,
    IN p_id_empleado INTEGER,
    IN p_id_cliente INTEGER,
    IN p_fecha_compra DATE,
    IN p_total_compra NUMERIC,
    IN p_id_producto INTEGER,
    IN p_cantidad_producto INTEGER
)
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE compra
    SET id_empleado = p_id_empleado,
        id_cliente = p_id_cliente,
        fecha_compra = p_fecha_compra,
        total_compra = p_total_compra
    WHERE id_compra = p_id_compra;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'compra_not_found';
    END IF;

    DELETE FROM producto_compra
    WHERE id_compra = p_id_compra;

    INSERT INTO producto_compra (
        id_compra,
        id_producto,
        cantidad_producto
    )
    VALUES (
        p_id_compra,
        p_id_producto,
        p_cantidad_producto
    );
END;
$$;
