// Bilingual option dictionaries — values must match PostgreSQL enums in migrations 003/004
export type Opt = { value: string; labelAm: string; labelEn: string };

export const OWNER_TYPES: Opt[] = [
  { value: "public", labelAm: "የሕዝብ", labelEn: "Public" },
  { value: "government", labelAm: "መንግሥት", labelEn: "Government" },
  { value: "religion", labelAm: "ሃይማኖት", labelEn: "Religious" },
  { value: "private", labelAm: "የግል", labelEn: "Private" },
  { value: "waqf", labelAm: "ወቂፍ", labelEn: "Waqf" },
];

export const AGE_METHODS: Opt[] = [
  { value: "estimated", labelAm: "ግምት", labelEn: "Estimated" },
  { value: "exact", labelAm: "Exact", labelEn: "Exact" },
  { value: "relative", labelAm: "Relative", labelEn: "Relative" },
];

export const OVERALL_CONDITIONS: Opt[] = [
  { value: "very_good", labelAm: "በጣም ጥሩ", labelEn: "Very Good" },
  { value: "good", labelAm: "ጥሩ", labelEn: "Good" },
  { value: "damaged", labelAm: "የተጎዳ", labelEn: "Damaged" },
  { value: "severely_damaged", labelAm: "በከፍተኛ ሁኔታ የተጎዳ", labelEn: "Severely Damaged" },
];

export const DAMAGE_LEVELS: Opt[] = [
  { value: "minor", labelAm: "አነስተኛ", labelEn: "Minor" },
  { value: "low", labelAm: "ዝቅተኛ", labelEn: "Low" },
  { value: "medium", labelAm: "መካከለኛ", labelEn: "Medium" },
  { value: "severe", labelAm: "ከፍተኛ", labelEn: "Severe" },
];

export const QUALITY_LEVELS: Opt[] = [
  { value: "very_good", labelAm: "በጣም ጥሩ", labelEn: "Very Good" },
  { value: "good", labelAm: "ጥሩ", labelEn: "Good" },
  { value: "medium", labelAm: "መካከለኛ", labelEn: "Medium" },
  { value: "low", labelAm: "ዝቅተኛ", labelEn: "Low" },
  { value: "very_low", labelAm: "በጣም ዝቅተኛ", labelEn: "Very Low" },
];

export const ACCESSIBILITY_LEVELS: Opt[] = [
  { value: "very_good", labelAm: "በጣም ጥሩ", labelEn: "Very Good" },
  { value: "good", labelAm: "ጥሩ", labelEn: "Good" },
  { value: "medium", labelAm: "መካከለኛ", labelEn: "Medium" },
  { value: "low", labelAm: "ዝቅተኛ", labelEn: "Low" },
  { value: "very_low", labelAm: "በጣም ዝቅተኛ", labelEn: "Very Low" },
  { value: "none", labelAm: "ምንም", labelEn: "None" },
];

export const SEX_TYPES: Opt[] = [
  { value: "male", labelAm: "ወንድ", labelEn: "Male" },
  { value: "female", labelAm: "ሴት", labelEn: "Female" },
];

// Heritage category codes — must match backend validCategories (I, II, VII, VIII)
export const CATEGORIES: Opt[] = [
  { value: "I", labelAm: "ምድብ I — ሃውልቶች", labelEn: "Category I — Monuments" },
  { value: "II", labelAm: "ምድብ II — ህንፃዎች", labelEn: "Category II — Buildings" },
  { value: "VII", labelAm: "ምድብ VII — የባህል እንጦጦ", labelEn: "Category VII — Cultural Landscapes" },
  { value: "VIII", labelAm: "ምድብ VIII — የታሪክ ትሁንባዊ ከተሞች", labelEn: "Category VIII — Historic Towns" },
];

export const CURRENT_USES: Opt[] = [
  { value: "residence", labelAm: "መኖሪያ", labelEn: "Residence" },
  { value: "worship", labelAm: "የአምልኮ ቦታ", labelEn: "Worship" },
  { value: "mosque", labelAm: "መስጂድ", labelEn: "Mosque" },
  { value: "museum", labelAm: "ሙዚየም", labelEn: "Museum" },
  { value: "shop", labelAm: "ሱቅ", labelEn: "Shop" },
  { value: "office", labelAm: "ቢሮ", labelEn: "Office" },
  { value: "abandoned", labelAm: "የተተወ", labelEn: "Abandoned" },
];

export const HARARI_HOUSE_GRADES: Opt[] = [
  { value: "dera_yalew", labelAm: "ዴራ ያለው", labelEn: "With Dera (inner room)" },
  { value: "dera_yeleew", labelAm: "ዴራ የሌለው", labelEn: "Without Dera" },
  { value: "bale_hulet_dera", labelAm: "ባለ ሁለት ዴራ", labelEn: "Double Dera" },
  { value: "bale_and_dera", labelAm: "ባለ አንድ ዴራ", labelEn: "Single Dera" },
  { value: "kirtet_yalew", labelAm: "ኪርተት ያለው", labelEn: "With Kirtet" },
  { value: "kirtet_yeleew", labelAm: "ኪርተት የሌለው", labelEn: "Without Kirtet" },
  { value: "qutiqele", labelAm: "ቁጢቀለ", labelEn: "Qutiqele (small upper)" },
  { value: "qela_gar", labelAm: "ቀላ ጋር", labelEn: "Qela Gar (upper floor)" },
  { value: "qela_yeleew", labelAm: "ቀላ የሌለው", labelEn: "Without Qela" },
  { value: "amirgar", labelAm: "አሚርጋር", labelEn: "Amirgar (royal house)" },
];

export const NEIGHBORHOOD_TYPES: Opt[] = [
  { value: "jugol", labelAm: "ጁጎል", labelEn: "Jugol (walled city)" },
  { value: "outside_jugol", labelAm: "ከጁጎል ውጭ", labelEn: "Outside Jugol" },
];

export const RELATED_DOCS: Opt[] = [
  { value: "book", labelAm: "መጽሐፍ", labelEn: "Book" },
  { value: "photo", labelAm: "ፎቶ ግራፍ", labelEn: "Photograph" },
  { value: "slide", labelAm: "ስላይድ", labelEn: "Slide" },
  { value: "map", labelAm: "ካርታ", labelEn: "Map" },
  { value: "register", labelAm: "መዝገብ", labelEn: "Register" },
  { value: "plan", labelAm: "ፕላን", labelEn: "Plan / Drawing" },
];
