CREATE TABLE positions
(
    ID uuid PRIMARY KEY ,
    name TEXT NOT NULL,
    salary INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()

)