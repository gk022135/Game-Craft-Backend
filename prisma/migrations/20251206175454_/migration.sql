/*
  Warnings:

  - You are about to drop the `Question_Records` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `Tables_Info` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropTable
DROP TABLE "public"."Question_Records";

-- DropTable
DROP TABLE "public"."Tables_Info";

-- CreateTable
CREATE TABLE "TablesInfo" (
    "id" SERIAL NOT NULL,
    "tableName" TEXT NOT NULL,
    "Description" TEXT NOT NULL,
    "querry" TEXT NOT NULL,
    "CreatedBy" TEXT NOT NULL,
    "CreatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "TablesInfo_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "QuestionRecords" (
    "id" SERIAL NOT NULL,
    "ContributedBy" TEXT NOT NULL,
    "Title" TEXT NOT NULL,
    "Description" TEXT NOT NULL,
    "Topics" TEXT[],
    "UsedTables" TEXT[],
    "CreatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "DifficultyLevel" TEXT NOT NULL,
    "Rewards" INTEGER,
    "Answer" TEXT,

    CONSTRAINT "QuestionRecords_pkey" PRIMARY KEY ("id")
);
