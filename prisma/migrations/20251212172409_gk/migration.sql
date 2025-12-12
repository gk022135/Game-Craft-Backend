-- CreateTable
CREATE TABLE "User" (
    "id" SERIAL NOT NULL,
    "FirstName" TEXT NOT NULL,
    "LastName" TEXT NOT NULL,
    "Email" TEXT NOT NULL,
    "Password" TEXT NOT NULL,
    "Username" TEXT NOT NULL,
    "SolvedEasy" INTEGER[] DEFAULT ARRAY[]::INTEGER[],
    "SolvedMedium" INTEGER[] DEFAULT ARRAY[]::INTEGER[],
    "SolvedHard" INTEGER[] DEFAULT ARRAY[]::INTEGER[],
    "EarnedPoints" INTEGER NOT NULL DEFAULT 0,
    "CurrentLevel" INTEGER NOT NULL DEFAULT 1,
    "CurrentStreak" INTEGER NOT NULL DEFAULT 0,
    "MaxStreak" INTEGER NOT NULL DEFAULT 0,
    "Links" TEXT[] DEFAULT ARRAY[]::TEXT[],
    "PhoneNumber" TEXT,
    "Badges" INTEGER[] DEFAULT ARRAY[]::INTEGER[],

    CONSTRAINT "User_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Otp" (
    "id" SERIAL NOT NULL,
    "Username" TEXT NOT NULL,
    "FirstName" TEXT NOT NULL,
    "LastName" TEXT NOT NULL,
    "Email" TEXT NOT NULL,
    "Otp" TEXT NOT NULL,
    "ExpiresAt" TIMESTAMP(3) NOT NULL,
    "Password" TEXT NOT NULL,

    CONSTRAINT "Otp_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Question" (
    "id" SERIAL NOT NULL,
    "Title" TEXT NOT NULL,
    "Description" TEXT NOT NULL,
    "StarterSchema" TEXT NOT NULL,
    "StarterData" TEXT NOT NULL,
    "CorrectQuery" TEXT NOT NULL,
    "EndingSchema" TEXT NOT NULL,

    CONSTRAINT "Question_pkey" PRIMARY KEY ("id")
);

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
    "TitleLowerCase" TEXT NOT NULL DEFAULT 'gaurav kumar',
    "Description" TEXT NOT NULL,
    "Topics" TEXT[],
    "UsedTables" TEXT[],
    "CreatedAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "DifficultyLevel" TEXT NOT NULL,
    "Rewards" INTEGER,
    "Answer" TEXT,

    CONSTRAINT "QuestionRecords_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "UserActivityLog" (
    "id" SERIAL NOT NULL,
    "userId" INTEGER NOT NULL,
    "questionId" INTEGER NOT NULL,
    "timestamp" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "solution" TEXT NOT NULL,
    "isValid" BOOLEAN NOT NULL DEFAULT false,

    CONSTRAINT "UserActivityLog_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "User_Email_key" ON "User"("Email");

-- CreateIndex
CREATE UNIQUE INDEX "User_Username_key" ON "User"("Username");

-- CreateIndex
CREATE UNIQUE INDEX "Otp_Username_key" ON "Otp"("Username");

-- CreateIndex
CREATE UNIQUE INDEX "Otp_Email_key" ON "Otp"("Email");

-- AddForeignKey
ALTER TABLE "UserActivityLog" ADD CONSTRAINT "UserActivityLog_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "UserActivityLog" ADD CONSTRAINT "UserActivityLog_questionId_fkey" FOREIGN KEY ("questionId") REFERENCES "QuestionRecords"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
