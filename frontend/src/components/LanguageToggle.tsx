'use client';

import { useI18n } from '@/lib/i18n';
import { Button } from '@/components/ui/button';

export function LanguageToggle() {
  const { language, setLanguage } = useI18n();

  const toggleLanguage = () => {
    setLanguage(language === 'en' ? 'fa' : 'en');
  };

  return (
    <Button
      variant="ghost"
      size="sm"
      onClick={toggleLanguage}
      className="w-9 h-9 font-bold transition-all duration-300 hover:bg-accent"
    >
      {language === 'en' ? 'FA' : 'EN'}
    </Button>
  );
}
