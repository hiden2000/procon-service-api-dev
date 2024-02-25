
-- MySQLデータベース初期化スクリプト

-- ユーザーテーブル (Users)
CREATE TABLE IF NOT EXISTS Users (
    UserID INT AUTO_INCREMENT PRIMARY KEY,
    Username VARCHAR(255) NOT NULL,
    Password VARCHAR(255) NOT NULL,
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    LastLogin TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE INDEX username_unique (Username)
);

-- 問題テーブル (Problems)
CREATE TABLE IF NOT EXISTS Problems (
    ProblemID INT AUTO_INCREMENT PRIMARY KEY,
    UserID INT NOT NULL,
    Title VARCHAR(255) NOT NULL,
    Description TEXT,
    Difficulty INT CHECK(Difficulty >= 1 AND Difficulty <= 5),
    CreatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (UserID) REFERENCES Users(UserID),
    INDEX difficulty_index (Difficulty),
    INDEX user_id_index (UserID)
);

-- 解答テーブル (Solutions)
CREATE TABLE IF NOT EXISTS Solutions (
    SolutionID INT AUTO_INCREMENT PRIMARY KEY,
    UserID INT NOT NULL,
    ProblemID INT NOT NULL,
    LanguageID INT NOT NULL,
    Code TEXT,
    SubmittedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (UserID) REFERENCES Users(UserID),
    FOREIGN KEY (ProblemID) REFERENCES Problems(ProblemID),
    INDEX user_id_index (UserID),
    INDEX problem_id_index (ProblemID),
    INDEX language_id_index (LanguageID)
);

-- 判定結果テーブル (ResultDetails)
CREATE TABLE IF NOT EXISTS ResultDetails (
    SolutionID INT PRIMARY KEY,
    TotalCases INT NOT NULL,
    CorrectCases INT NOT NULL,
    IncorrectCases INT NOT NULL,
    TimeLimitExceeded INT NOT NULL,
    ErrorMessage TEXT,
    FOREIGN KEY (SolutionID) REFERENCES Solutions(SolutionID)
);

-- ケースごとの結果テーブル (CaseResults)
CREATE TABLE IF NOT EXISTS CaseResults (
    SolutionID INT,
    CaseName VARCHAR(255) NOT NULL,
    Result VARCHAR(255) NOT NULL,
    ExecutionTime INT NOT NULL,
    PRIMARY KEY (SolutionID, CaseName),
    FOREIGN KEY (SolutionID) REFERENCES Solutions(SolutionID)
);
