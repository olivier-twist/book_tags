create TABLE BOOK (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255),
    authorLF VARCHAR(255),
    additionalAuthors VARCHAR(255) ,
    tag VARCHAR(255) NOT NULL
);