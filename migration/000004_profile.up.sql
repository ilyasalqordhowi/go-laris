CREATE TABLE "profile"(
    "id" serial primary key,
    "fullname" VARCHAR(255),
    "province" VARCHAR(255),
    "city" VARCHAR(255),
    "postal_code" INTEGER,
    "country" VARCHAR(50),
    "mobile" VARCHAR(50),
    "address" VARCHAR(255)
);