CREATE TABLE climb (
    Id SERIAL PRIMARY KEY,
    Name VARCHAR(255),
    Grade VARCHAR(255),
    CragID INTEGER REFERENCES crag(Id)
);
-- name: CreateClimb :one 
-- columns: Id, Name, Grade, CragID
INSERT INTO climb (
    Name, Grade, CragID)
VALUES($1, $2, $3)
RETURNING *;

-- name: getClimbsByCrag :many
SELECT * FROM climb WHERE CragID = $1;

-- name: GetAllClimbs :many
SELECT * FROM climb ORDER BY name;

-- name: GetClimb :one
SELECT * FROM climb WHERE Id = $1;

-- name: UpdateClimb :exec
UPDATE climb set Name = $1, Grade = $2, CragID = $3 WHERE Id = $4;

-- name: DeleteClimb :exec
DELETE FROM climb WHERE id = $1;