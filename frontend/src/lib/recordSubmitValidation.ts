import type { ImmovableRecordInput, MovableRecordInput } from "@/types";

const IMMOVABLE_CATEGORIES = new Set(["I", "II", "VII", "VIII"]);
const MOVABLE_CATEGORIES = new Set(["III", "IV", "V", "VI"]);

export type SubmitFieldKey =
  | "name_amharic"
  | "category"
  | "woreda"
  | "kebele"
  | "current_use"
  | "owner_type"
  | "construction_period"
  | "age_method"
  | "storage_location"
  | "materials";

function hasValidImmovableCategory(categories: string[] | undefined): boolean {
  return (categories ?? []).some((c) => IMMOVABLE_CATEGORIES.has(c.trim().toUpperCase()));
}

function hasValidMovableCategory(category: string | undefined): boolean {
  if (!category?.trim()) return false;
  return MOVABLE_CATEGORIES.has(category.trim().toUpperCase());
}

export function validateImmovableForSubmit(form: ImmovableRecordInput): SubmitFieldKey[] {
  const missing: SubmitFieldKey[] = [];

  if (!form.name_amharic?.trim()) missing.push("name_amharic");
  if (!hasValidImmovableCategory(form.category)) missing.push("category");
  if (!form.woreda?.trim()) missing.push("woreda");
  if (!form.kebele?.trim()) missing.push("kebele");
  if (!form.current_use?.length) missing.push("current_use");
  if (!form.owner_type) missing.push("owner_type");
  if (!form.construction_period?.trim()) missing.push("construction_period");
  if (!form.age_method) missing.push("age_method");

  return missing;
}

export function validateMovableForSubmit(form: MovableRecordInput): SubmitFieldKey[] {
  const missing: SubmitFieldKey[] = [];

  if (!form.name_amharic?.trim()) missing.push("name_amharic");
  if (!hasValidMovableCategory(form.category)) missing.push("category");
  if (!form.woreda?.trim()) missing.push("woreda");
  if (!form.kebele?.trim()) missing.push("kebele");
  if (!form.owner_type) missing.push("owner_type");
  if (!form.storage_location) missing.push("storage_location");
  if (!form.materials?.length) missing.push("materials");

  return missing;
}
