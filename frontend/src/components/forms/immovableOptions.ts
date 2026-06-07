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
  { value: "minor", labelAm: "ጥቂት", labelEn: "Minor" },
  { value: "moderate", labelAm: "መካከለኛ", labelEn: "Moderate" },
  { value: "medium", labelAm: "Medium", labelEn: "Medium" },
  { value: "severe", labelAm: "ከባድ", labelEn: "Severe" },
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
  { value: "museum", labelAm: "ሙዚየም", labelEn: "Museum" },
  { value: "shop", labelAm: "ሱቅ", labelEn: "Shop" },
  { value: "office", labelAm: "ቢሮ", labelEn: "Office" },
  { value: "abandoned", labelAm: "የተተወ", labelEn: "Abandoned" },
];

export const HARARI_HOUSE_GRADES: Opt[] = [
  { value: "gar_gar", labelAm: "ጋር ጋር", labelEn: "Gar Gar" },
  { value: "harari_gar", labelAm: "ሐረሪ ጋር", labelEn: "Harari Gar" },
  { value: "indian_style", labelAm: "ህንዳዊ ቅርጽ", labelEn: "Indian Style" },
  { value: "modern", labelAm: "ዘመናዊ", labelEn: "Modern" },
];

export const NEIGHBORHOOD_TYPES: Opt[] = [
  { value: "jugol", labelAm: "ጁጎል", labelEn: "Jugol (walled city)" },
  { value: "outside_jugol", labelAm: "ከጁጎል ውጭ", labelEn: "Outside Jugol" },
];
