-- CreateTable
CREATE TABLE "Tables_Info" (
    "id" SERIAL NOT NULL,
    "tableName" TEXT NOT NULL,
    "Description" TEXT NOT NULL,
    "querry" TEXT NOT NULL,
    "CreatedBy" TEXT NOT NULL,
    "CreatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "Tables_Info_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Question_Records" (
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

    CONSTRAINT "Question_Records_pkey" PRIMARY KEY ("id")
);
