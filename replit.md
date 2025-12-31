# Admin Panel Project

## Overview
A modern, production-grade Admin Panel built with Next.js (Frontend) and Go (Backend). It features a robust theme system, internationalization (English/Persian), and a clean, high-performance UI.

## Recent Changes (December 31, 2025)
- **Theme System:** Integrated `next-themes` with support for Light, Dark, and a custom "Legendary" theme.
- **UI Modernization:** Refactored the Dashboard with entrance animations, backdrop-blur effects, and improved card designs.
- **Theme/Language Toggle:** Added a unified toggle system in the header for switching between themes and languages (EN/FA).
- **Icons:** Switched to Lucide React for consistent and modern iconography.

## Project Architecture
- **Frontend:** Next.js 16+, Tailwind CSS, TanStack Query, `next-themes`, Lucide React.
- **Backend:** Go, Chi router, PostgreSQL.
- **State Management:** TanStack Query for server state, custom I18n provider for client state.

## Authentication Contract
- **Login/Refresh Response Shape:** Access token and `token_expires_at` are returned in JSON responses. The `refresh_token` is issued only as an HttpOnly cookie (no JSON field for it).
- **Frontend Storage Rules:** The frontend may store only the access token and `token_expires_at`. Storing the refresh token in `localStorage` or any JS-accessible storage is forbidden.
- **Migration Note:** Remove any legacy `refresh_token` value from `localStorage`. Users may need to sign in again after this change.
- **Rollout/Rollback Notes:** Roll out by enabling the HttpOnly refresh cookie and disabling refresh token persistence on the client. If rolling back, be aware that users might have cleared localStorage refresh tokens, so expect another login prompt; verify refresh endpoints accept cookies and ensure fallback messaging.

## User Preferences
- Clean, minimalist design with attention to spacing and typography.
- Support for RTL (Persian) and LTR (English) layouts.
- Performance-oriented component design.
