'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';

export type Language = 'en' | 'fa';

interface Translations {
  [key: string]: {
    en: string;
    fa: string;
  };
}

const translations: Translations = {
  'app.title': { en: 'Admin Panel', fa: 'پنل مدیریت' },
  'app.signIn': { en: 'Sign in to your account', fa: 'وارد حساب کاربری شوید' },
  'nav.dashboard': { en: 'Dashboard', fa: 'داشبورد' },
  'nav.users': { en: 'Users', fa: 'کاربران' },
  'nav.roles': { en: 'Roles', fa: 'نقش‌ها' },
  'nav.auditLogs': { en: 'Audit Logs', fa: 'گزارش‌ها' },
  'nav.settings': { en: 'Settings', fa: 'تنظیمات' },
  'nav.featureFlags': { en: 'Feature Flags', fa: 'ویژگی‌ها' },
  'auth.email': { en: 'Email', fa: 'ایمیل' },
  'auth.password': { en: 'Password', fa: 'رمز عبور' },
  'auth.signIn': { en: 'Sign in', fa: 'ورود' },
  'auth.logout': { en: 'Logout', fa: 'خروج' },
  'auth.demo': { en: 'Demo:', fa: 'نسخه آزمایشی:' },
  'theme.title': { en: 'Theme', fa: 'تم' },
  'theme.light': { en: 'Light', fa: 'روشن' },
  'theme.dark': { en: 'Dark', fa: 'تاریک' },
  'theme.legendary': { en: 'Legendary', fa: 'افسانه‌ای' },
  'language.title': { en: 'Language', fa: 'زبان' },
  'dashboard.title': { en: 'Dashboard', fa: 'داشبورد' },
  'dashboard.welcome': { en: 'Welcome back', fa: 'خوش آمدید' },
  'dashboard.totalUsers': { en: 'Total Users', fa: 'کل کاربران' },
  'dashboard.activeUsers': { en: 'Active Users', fa: 'کاربران فعال' },
  'dashboard.totalRoles': { en: 'Total Roles', fa: 'کل نقش‌ها' },
  'dashboard.auditLogs': { en: 'Audit Logs', fa: 'گزارش‌ها' },
  'users.title': { en: 'User Management', fa: 'مدیریت کاربران' },
  'users.search': { en: 'Search users...', fa: 'جستجوی کاربران...' },
  'users.addUser': { en: 'Add User', fa: 'افزودن کاربر' },
  'users.name': { en: 'Name', fa: 'نام' },
  'users.email': { en: 'Email', fa: 'ایمیل' },
  'users.role': { en: 'Role', fa: 'نقش' },
  'users.status': { en: 'Status', fa: 'وضعیت' },
  'users.actions': { en: 'Actions', fa: 'عملیات' },
  'users.active': { en: 'Active', fa: 'فعال' },
  'users.inactive': { en: 'Inactive', fa: 'غیرفعال' },
  'roles.title': { en: 'Role Management', fa: 'مدیریت نقش‌ها' },
  'roles.addRole': { en: 'Add Role', fa: 'افزودن نقش' },
  'audit.title': { en: 'Audit Logs', fa: 'گزارش‌ها' },
  'audit.export': { en: 'Export CSV', fa: 'خروجی CSV' },
  'settings.title': { en: 'Settings', fa: 'تنظیمات' },
  'common.loading': { en: 'Loading...', fa: 'در حال بارگذاری...' },
  'common.save': { en: 'Save', fa: 'ذخیره' },
  'common.cancel': { en: 'Cancel', fa: 'انصراف' },
  'common.delete': { en: 'Delete', fa: 'حذف' },
  'common.edit': { en: 'Edit', fa: 'ویرایش' },
  'common.create': { en: 'Create', fa: 'ایجاد' },
  'common.search': { en: 'Search', fa: 'جستجو' },
  'common.noResults': { en: 'No results found', fa: 'نتیجه‌ای یافت نشد' },
  'common.previous': { en: 'Previous', fa: 'قبلی' },
  'common.next': { en: 'Next', fa: 'بعدی' },
};

interface I18nContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string) => string;
  isRTL: boolean;
  mounted: boolean;
}

const I18nContext = createContext<I18nContextType | undefined>(undefined);

export function I18nProvider({ children }: { children: ReactNode }) {
  const [language, setLanguageState] = useState<Language>('en');
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    const stored = localStorage.getItem('language') as Language | null;
    if (stored && (stored === 'en' || stored === 'fa')) {
      setLanguageState(stored);
    }
    setMounted(true);
  }, []);

  const setLanguage = (lang: Language) => {
    setLanguageState(lang);
    localStorage.setItem('language', lang);
    document.documentElement.dir = lang === 'fa' ? 'rtl' : 'ltr';
    document.documentElement.lang = lang;
  };

  useEffect(() => {
    if (mounted) {
      document.documentElement.dir = language === 'fa' ? 'rtl' : 'ltr';
      document.documentElement.lang = language;
    }
  }, [language, mounted]);

  const t = (key: string): string => {
    const translation = translations[key];
    if (!translation) return key;
    return translation[language] || translation.en || key;
  };

  const isRTL = language === 'fa';

  return (
    <I18nContext.Provider value={{ language, setLanguage, t, isRTL, mounted }}>
      {children}
    </I18nContext.Provider>
  );
}

export function useI18n() {
  const context = useContext(I18nContext);
  if (context === undefined) {
    throw new Error('useI18n must be used within an I18nProvider');
  }
  return context;
}
