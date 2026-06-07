import type { Opt } from "./immovableOptions";

// Heritage category codes — must match backend validCategories (III–VI)
export const MOVABLE_CATEGORIES: Opt[] = [
  { value: "III", labelAm: "ምድብ III — የቅርስ ነገሮች", labelEn: "Category III — Movable Heritage" },
  { value: "IV", labelAm: "ምድብ IV — ታሪክያዊ ፍጻሜዎች", labelEn: "Category IV — Historical Events" },
  { value: "V", labelAm: "ምድብ V — ባህላዊ ማረፊያዎች", labelEn: "Category V — Traditional Settlements" },
  { value: "VI", labelAm: "ምድብ VI — የአምልኮ ቦታዎች", labelEn: "Category VI — Cult Sites" },
];

export const MOVABLE_OWNER_TYPES: Opt[] = [
  { value: "public", labelAm: "የሕዝብ", labelEn: "Public" },
  { value: "government", labelAm: "መንግሥት", labelEn: "Government" },
  { value: "religion", labelAm: "ሃይማኖት", labelEn: "Religious" },
  { value: "private", labelAm: "የግል", labelEn: "Private" },
];

export const STORAGE_LOCATIONS: Opt[] = [
  { value: "museum", labelAm: "ሙዚየም", labelEn: "Museum" },
  { value: "store", labelAm: "መጋዘን", labelEn: "Store" },
  { value: "church", labelAm: "ቤተ ክርስቲያን", labelEn: "Church" },
  { value: "private_home", labelAm: "የግል ቤት", labelEn: "Private Home" },
  { value: "other", labelAm: "ሌላ", labelEn: "Other" },
];

export const MOVABLE_AGE_METHODS: Opt[] = [
  { value: "estimated", labelAm: "ግምት", labelEn: "Estimated" },
  { value: "exact", labelAm: "Exact", labelEn: "Exact" },
  { value: "relative", labelAm: "Relative", labelEn: "Relative" },
];

export const MOVABLE_CONDITIONS: Opt[] = [
  { value: "good", labelAm: "ጥሩ", labelEn: "Good" },
  { value: "fair", labelAm: "መካከለኛ", labelEn: "Fair" },
  { value: "damaged", labelAm: "የተጎዳ", labelEn: "Damaged" },
  { value: "incomplete", labelAm: "Incomplete", labelEn: "Incomplete" },
];

export const ACQUISITION_METHODS: Opt[] = [
  { value: "purchase", labelAm: "ግዢ", labelEn: "Purchase" },
  { value: "donation", labelAm: "ስጦታ", labelEn: "Donation" },
  { value: "inheritance", labelAm: "ውርስ", labelEn: "Inheritance" },
  { value: "transfer", labelAm: "ዝውውር", labelEn: "Transfer" },
  { value: "found", labelAm: "የተገኘ", labelEn: "Found" },
];

export const MATERIALS: Opt[] = [
  { value: "wood", labelAm: "እንጨት", labelEn: "Wood" },
  { value: "metal", labelAm: "ብረት", labelEn: "Metal" },
  { value: "textile", labelAm: "ጨርቅ", labelEn: "Textile" },
  { value: "leather", labelAm: "ቆዳ", labelEn: "Leather" },
  { value: "stone", labelAm: "ድንጋይ", labelEn: "Stone" },
  { value: "paper", labelAm: "ወረቀት", labelEn: "Paper" },
  { value: "ceramic", labelAm: "ሴራሚክ", labelEn: "Ceramic" },
  { value: "other", labelAm: "ሌላ", labelEn: "Other" },
];

export const NOTABLE_BECAUSE: Opt[] = [
  { value: "age", labelAm: "ዕድሜ", labelEn: "Age" },
  { value: "rarity", labelAm: "አልፎ አልፎ", labelEn: "Rarity" },
  { value: "craftsmanship", labelAm: "የእጅ ሥራ", labelEn: "Craftsmanship" },
  { value: "historical", labelAm: "ታሪካዊ", labelEn: "Historical" },
  { value: "religious", labelAm: "ሃይማኖታዊ", labelEn: "Religious" },
  { value: "cultural", labelAm: "ባህላዊ", labelEn: "Cultural" },
];

export {
  QUALITY_LEVELS,
  ACCESSIBILITY_LEVELS,
  SEX_TYPES,
  AGE_METHODS,
} from "./immovableOptions";
