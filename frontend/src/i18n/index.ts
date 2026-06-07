import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import am from "./am.json";
import en from "./en.json";

if (!i18n.isInitialized) {
  void i18n.use(initReactI18next).init({
    resources: {
      am: { translation: am },
      en: { translation: en },
    },
    lng: "am",
    fallbackLng: "en",
    interpolation: { escapeValue: false },
    react: { useSuspense: false },
  });
}

export default i18n;
