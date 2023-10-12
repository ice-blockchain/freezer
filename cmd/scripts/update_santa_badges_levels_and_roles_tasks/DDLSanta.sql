-- SPDX-License-Identifier: ice License 1.0
--************************************************************************************************************************************
-- badge_progress
CREATE TABLE IF NOT EXISTS badge_progress (
                        balance          BIGINT NOT NULL DEFAULT 0,
                        friends_invited  BIGINT NOT NULL DEFAULT 0,
                        completed_levels BIGINT NOT NULL DEFAULT 0,
                        hide_badges      BOOLEAN DEFAULT FALSE,
                        achieved_badges  TEXT[],
                        user_id          TEXT NOT NULL PRIMARY KEY
                    ) WITH (fillfactor = 70);
--************************************************************************************************************************************
-- badge_statistics
CREATE TABLE IF NOT EXISTS badge_statistics (
                        achieved_by        BIGINT NOT NULL DEFAULT 0,
                        badge_type         TEXT NOT NULL PRIMARY KEY,
                        badge_group_type   TEXT NOT NULL
                    ) WITH (fillfactor = 70);

-- levels_and_roles_progress
CREATE TABLE IF NOT EXISTS levels_and_roles_progress (
                        mining_streak               BIGINT NOT NULL DEFAULT 0,
                        pings_sent                  BIGINT NOT NULL DEFAULT 0,
                        friends_invited             BIGINT NOT NULL DEFAULT 0,
                        completed_tasks             BIGINT NOT NULL DEFAULT 0,
                        hide_level                  BOOLEAN DEFAULT false,
                        hide_role                   BOOLEAN DEFAULT false,
                        agenda_contact_user_ids     TEXT[],
                        enabled_roles               TEXT[],
                        completed_levels            TEXT[],
                        user_id                     TEXT NOT NULL PRIMARY KEY,
                        phone_number_hash           TEXT
                    ) WITH (fillfactor = 70);

-- task_progress
CREATE TABLE IF NOT EXISTS task_progress (
                        friends_invited             BIGINT NOT NULL DEFAULT 0,
                        mining_started              BOOLEAN DEFAULT FALSE,
                        username_set                BOOLEAN DEFAULT FALSE,
                        profile_picture_set         BOOLEAN DEFAULT FALSE,
                        completed_tasks             TEXT[],
                        pseudo_completed_tasks      TEXT[],
                        user_id                     TEXT NOT NULL PRIMARY KEY,
                        twitter_user_handle         TEXT,
                        telegram_user_handle        TEXT
                    ) WITH (fillfactor = 70);