━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📄 SPEC BLOCK — Phase 1: Project Vision
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PROJECT VISION
==============
Project Name:    ቅርስ መዝገብ (Qirs Mezgeb) — Heritage Registry System
Type:            Government Management System (Web + Mobile PWA)
Core Problem:    Heritage registration in the Harari region is done
                 manually on paper forms, making records hard to search,
                 preserve, share, and approve through a chain of authority.
Who Suffers:     Field registrars (data entry in the field), supervisors
                 (no visibility on submission status), managers (no central
                 oversight or audit trail), and the institution itself
                 (records get lost, damaged, or duplicated).
Why Current
Solutions Fail:  Paper forms have no approval workflow, no search, no
                 photo attachment, no status tracking, and no protection
                 against loss or duplication.
What This
System Improves: A structured digital workflow — from field data entry to
                 final approval — with a searchable, permanent, bilingual
                 record for every heritage asset.
Success Looks
Like:            A registrar fills a form on their phone in the field,
                 a supervisor reviews it the same day, a manager approves
                 it with one click, and the record is permanently stored
                 with a unique ID, searchable by anyone in the institution.
                 ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📄 SPEC BLOCK — Phase 2: System Actors
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

SYSTEM ACTORS & PERMISSIONS
============================

Actor: Registrar
  Can View:    Own submitted records, own drafts, returned records + comments
  Can Create:  New immovable heritage record, new movable heritage record
  Can Update:  Own draft records, returned records (after supervisor feedback)
  Can Delete:  Own drafts only (not submitted records)
  Depends On:  Supervisor to review, Manager for final approval
  Workflow:    Log in → choose form type → fill form → save draft or submit
               → receive notification if returned → fix and resubmit

Actor: Supervisor
  Can View:    All submitted records (any registrar), all reviewed records
  Can Create:  Review comments on any record
  Can Update:  Record status (approve to manager / return to registrar)
  Can Delete:  Nothing
  Depends On:  Registrar to submit, Manager for final decision
  Workflow:    Log in → see pending submissions → open record → review
               → approve (sends to manager) or return with comments

Actor: Manager
  Can View:    All records at any status, full dashboard, all reports
  Can Create:  New user accounts, role assignments
  Can Update:  Record status (final approval), user roles
  Can Delete:  User accounts only
  Depends On:  Supervisor to pre-review before records reach them
  Workflow:    Log in → see supervisor-approved records → final approve
               → manage users → view reports & export records

               ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📄 SPEC BLOCK — Phase 3: Module Breakdown
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

MODULE BREAKDOWN
================

Module: Auth & User Management
  Primary Actor(s): Manager (admin), all actors (login)
  Core Responsibility: Secure login, role assignment, session control
  Sub-features: Login, logout, JWT token, role-based route guards,
                create/edit/deactivate users, language preference (AM/EN)
  Depends On: Nothing — must be built first
  Priority: 1 (must-have)

Module: Immovable Heritage Registration (Form 02)
  Primary Actor(s): Registrar
  Core Responsibility: Digital Form 02 — full data entry for immovable assets
  Sub-features: All form sections, photo upload, GPS capture,
                save as draft, submit for review, auto-generate ID (ET-HR-AN-I-XXXX)
  Depends On: Auth module
  Priority: 1 (must-have)

Module: Movable Heritage Registration (Form 01)
  Primary Actor(s): Registrar
  Core Responsibility: Digital Form 01 — full data entry for movable assets
  Sub-features: All form sections, photo upload, measurements,
                material checkboxes, save as draft, submit,
                auto-generate ID (ET-HR-AN-V-XXXX)
  Depends On: Auth module
  Priority: 1 (must-have)

Module: Approval Workflow
  Primary Actor(s): Supervisor, Manager
  Core Responsibility: 3-stage status pipeline with comments and notifications
  Sub-features: Pending queue per role, approve/return action,
                comment thread per record, status history/audit log,
                in-app notification badge
  Depends On: Both registration modules
  Priority: 1 (must-have)

Module: Dashboard & Search
  Primary Actor(s): All actors (filtered by role)
  Core Responsibility: Overview stats, record search, filters
  Sub-features: Record count by status/type/region, search by name/ID/location,
                filter by form type / status / date, export to PDF/CSV
  Depends On: Both registration modules + approval workflow
  Priority: 2 (important)

Module: Reports & Export
  Primary Actor(s): Manager, Supervisor
  Core Responsibility: Printable record sheets and data exports
  Sub-features: Print single record as formatted PDF,
                export filtered list to CSV, summary statistics report
  Depends On: Dashboard module
  Priority: 2 (important)