import { useState, type FormEvent } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Save, Send, Loader2 } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type {
  ImmovableRecord,
  ImmovableRecordInput,
  RecordPhoto,
} from "@/types";
import {
  Field,
  FormSection,
  TextInput,
  NumberInput,
  TextArea,
  Select,
  CheckboxList,
  Switch,
} from "./Field";
import { PhotoUploader } from "./PhotoUploader";
import {
  createImmovable,
  updateImmovable,
  submitImmovable,
  uploadImmovablePhoto,
  deleteImmovablePhoto,
} from "@/api/immovable";
import { toDateInputValue } from "@/lib/recordPayload";
import { validateImmovableForSubmit } from "@/lib/recordSubmitValidation";
import { getApiErrorMessage } from "@/lib/apiError";
import {
  OWNER_TYPES,
  AGE_METHODS,
  OVERALL_CONDITIONS,
  DAMAGE_LEVELS,
  QUALITY_LEVELS,
  ACCESSIBILITY_LEVELS,
  SEX_TYPES,
  CATEGORIES,
  CURRENT_USES,
  HARARI_HOUSE_GRADES,
  NEIGHBORHOOD_TYPES,
  RELATED_DOCS,
} from "./immovableOptions";

interface ImmovableFormProps {
  mode: "create" | "edit";
  initialRecord?: ImmovableRecord;
  photos?: RecordPhoto[];
}

// Helper to coerce empty string -> undefined for optional numbers
const num = (v: string) => (v === "" ? undefined : Number(v));

export function ImmovableForm({ mode, initialRecord, photos = [] }: ImmovableFormProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const qc = useQueryClient();

  const [form, setForm] = useState<ImmovableRecordInput>(() => ({
    name_amharic: initialRecord?.name_amharic ?? "",
    name_local: initialRecord?.name_local ?? "",
    category: initialRecord?.category ?? [],
    current_use: initialRecord?.current_use ?? [],
    current_use_other: initialRecord?.current_use_other ?? "",
    previous_id: initialRecord?.previous_id ?? "",
    woreda: initialRecord?.woreda ?? "",
    kebele: initialRecord?.kebele ?? "",
    house_number: initialRecord?.house_number ?? "",
    street_number: initialRecord?.street_number ?? "",
    gate: initialRecord?.gate ?? "",
    owner_type: initialRecord?.owner_type,
    owner_name: initialRecord?.owner_name ?? "",
    map_reference: initialRecord?.map_reference ?? "",
    gps_east: initialRecord?.gps_east,
    gps_north: initialRecord?.gps_north,
    elevation_m: initialRecord?.elevation_m,
    built_by: initialRecord?.built_by ?? "",
    construction_period: initialRecord?.construction_period ?? "",
    age_method: initialRecord?.age_method,
    height_m: initialRecord?.height_m,
    length_m: initialRecord?.length_m,
    width_m: initialRecord?.width_m,
    num_doors: initialRecord?.num_doors,
    num_windows: initialRecord?.num_windows,
    num_rooms: initialRecord?.num_rooms,
    material: initialRecord?.material ?? "",
    harari_house_grades: initialRecord?.harari_house_grades ?? [],
    neighborhood_type: initialRecord?.neighborhood_type ?? "",
    description: initialRecord?.description ?? "",
    overall_condition: initialRecord?.overall_condition,
    damage_roof: initialRecord?.damage_roof,
    damage_cornice: initialRecord?.damage_cornice,
    damage_wall: initialRecord?.damage_wall,
    damage_floor: initialRecord?.damage_floor,
    damage_door: initialRecord?.damage_door,
    damage_cupboard: initialRecord?.damage_cupboard,
    damage_upper_floor: initialRecord?.damage_upper_floor,
    damage_dera: initialRecord?.damage_dera,
    damage_pillar: initialRecord?.damage_pillar,
    value_historical: initialRecord?.value_historical ?? "",
    value_craftsmanship: initialRecord?.value_craftsmanship ?? "",
    value_artistic: initialRecord?.value_artistic ?? "",
    value_scientific: initialRecord?.value_scientific ?? "",
    value_cultural: initialRecord?.value_cultural ?? "",
    has_threat: initialRecord?.has_threat ?? false,
    maintenance_done: initialRecord?.maintenance_done ?? false,
    maintenance_reason: initialRecord?.maintenance_reason ?? "",
    maintenance_by: initialRecord?.maintenance_by ?? "",
    maintenance_date: toDateInputValue(initialRecord?.maintenance_date),
    maintenance_count: initialRecord?.maintenance_count,
    preventive_level: initialRecord?.preventive_level,
    accessibility: initialRecord?.accessibility,
    notes: initialRecord?.notes ?? "",
    has_oral_history: initialRecord?.has_oral_history ?? false,
    caretaker_name: initialRecord?.caretaker_name ?? "",
    caretaker_role: initialRecord?.caretaker_role ?? "",
    informant_name: initialRecord?.informant_name ?? "",
    informant_sex: initialRecord?.informant_sex,
    informant_age: initialRecord?.informant_age,
    registrar_date: toDateInputValue(initialRecord?.registrar_date),
  }));

  const [openSections, setOpenSections] = useState<Record<number, boolean>>({
    1: true, 2: true, 3: false, 4: false, 5: false, 6: false, 7: false, 8: false,
  });
  const toggle = (n: number) =>
    setOpenSections((s) => ({ ...s, [n]: !s[n] }));

  const [error, setError] = useState<string | null>(null);

  const update = <K extends keyof ImmovableRecordInput>(
    key: K,
    value: ImmovableRecordInput[K],
  ) => setForm((f) => ({ ...f, [key]: value }));

  const saveMut = useMutation({
    mutationFn: async (body: ImmovableRecordInput) => {
      if (mode === "edit" && initialRecord) {
        await updateImmovable(initialRecord.id, body);
        return initialRecord.id;
      }
      const result = await createImmovable(body);
      return result.id;
    },
    onSuccess: (id) => {
      qc.invalidateQueries({ queryKey: ["records"] });
      qc.invalidateQueries({ queryKey: ["immovable", id] });
      qc.invalidateQueries({ queryKey: ["record-detail", "immovable", id] });
    },
    onError: (e: unknown) => setError(getApiErrorMessage(e)),
  });

  const submitMut = useMutation({
    mutationFn: async () => {
      const id =
        mode === "edit" && initialRecord
          ? initialRecord.id
          : await saveMut.mutateAsync(form);
      // Persist current edits before submit
      if (mode === "edit" && initialRecord) {
        await updateImmovable(initialRecord.id, form);
      }
      await submitImmovable(id);
      return id;
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["records"] });
      qc.invalidateQueries({ queryKey: ["dashboard"] });
      navigate({ to: "/records" });
    },
    onError: (e: unknown) => setError(getApiErrorMessage(e)),
  });

  const handleSave = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    if (!form.name_amharic?.trim()) {
      setError(t("forms.required.nameAmharic"));
      return;
    }
    if (!form.woreda?.trim() || !form.kebele?.trim()) {
      setError(t("forms.required.location"));
      return;
    }
    const id = await saveMut.mutateAsync(form).catch(() => null);
    if (!id) return;
    if (mode === "create") navigate({ to: "/records/$id/edit", params: { id } });
  };

  const handleSubmit = () => {
    setError(null);
    const missing = validateImmovableForSubmit(form);
    if (missing.length > 0) {
      setOpenSections((s) => ({
        ...s,
        1:
          s[1] ||
          missing.some((f) =>
            ["name_amharic", "category", "current_use"].includes(f),
          ),
        2: s[2] || missing.some((f) => ["woreda", "kebele"].includes(f)),
        3:
          s[3] ||
          missing.some((f) =>
            ["owner_type", "construction_period", "age_method"].includes(f),
          ),
      }));
      const labels = missing.map((f) => t(`forms.fields.${f}`)).join(", ");
      setError(t("forms.submitMissing", { fields: labels }));
      return;
    }
    submitMut.mutate();
  };

  return (
    <form onSubmit={handleSave} className="space-y-4 pb-24">
      <div className="rounded-xl border border-primary/20 bg-primary/5 p-4">
        <h2 className="font-amharic text-lg font-bold text-foreground">
          {t("immovable.title")}
        </h2>
        <p className="font-amharic mt-1 text-xs text-muted-foreground">
          {t("immovable.subtitle")}
        </p>
        {initialRecord && (
          <p className="mt-2 text-xs text-muted-foreground tabular-nums">
            ID: {initialRecord.record_id} · {t(`status.${initialRecord.status}`)}
          </p>
        )}
      </div>

      {/* Section 1: Identification */}
      <FormSection num={1} am="መለያ" en="Identification" open={openSections[1]} onToggle={() => toggle(1)}>
        <Field am="ስም (አማርኛ)" en="Name (Amharic)" required span={2}>
          <TextInput
            value={form.name_amharic ?? ""}
            onChange={(e) => update("name_amharic", e.target.value)}
            required
          />
        </Field>
        <Field am="ሌላ ስም" en="Local / Other Name">
          <TextInput
            value={form.name_local ?? ""}
            onChange={(e) => update("name_local", e.target.value)}
          />
        </Field>
        <Field am="ምድብ" en="Category" required span={3}>
          <CheckboxList
            options={CATEGORIES}
            value={form.category ?? []}
            onChange={(v) => update("category", v)}
          />
        </Field>
        <Field am="የአሁን አጠቃቀም" en="Current Use" required span={2}>
          <CheckboxList
            options={CURRENT_USES}
            value={form.current_use ?? []}
            onChange={(v) => update("current_use", v)}
          />
        </Field>
        <Field am="ሌላ (ይግለጹ)" en="Other use">
          <TextInput
            value={form.current_use_other ?? ""}
            onChange={(e) => update("current_use_other", e.target.value)}
          />
        </Field>
        <Field am="የቀድሞ መለያ" en="Previous ID">
          <TextInput
            value={form.previous_id ?? ""}
            onChange={(e) => update("previous_id", e.target.value)}
          />
        </Field>
      </FormSection>

      {/* Section 2: Location */}
      <FormSection num={2} am="አድራሻ" en="Location" open={openSections[2]} onToggle={() => toggle(2)}>
        <Field am="ወረዳ" en="Woreda" required>
          <TextInput value={form.woreda ?? ""} onChange={(e) => update("woreda", e.target.value)} required />
        </Field>
        <Field am="ቀበሌ" en="Kebele" required>
          <TextInput value={form.kebele ?? ""} onChange={(e) => update("kebele", e.target.value)} required />
        </Field>
        <Field am="የቤት ቁጥር" en="House Number">
          <TextInput value={form.house_number ?? ""} onChange={(e) => update("house_number", e.target.value)} />
        </Field>
        <Field am="የመንገድ ቁጥር" en="Street Number">
          <TextInput value={form.street_number ?? ""} onChange={(e) => update("street_number", e.target.value)} />
        </Field>
        <Field am="በር" en="Gate">
          <TextInput value={form.gate ?? ""} onChange={(e) => update("gate", e.target.value)} />
        </Field>
        <Field am="የካርታ ማጣቀሻ" en="Map Reference">
          <TextInput value={form.map_reference ?? ""} onChange={(e) => update("map_reference", e.target.value)} />
        </Field>
        <Field am="GPS ምስራቅ" en="GPS East">
          <NumberInput value={form.gps_east ?? ""} onChange={(e) => update("gps_east", num(e.target.value))} />
        </Field>
        <Field am="GPS ሰሜን" en="GPS North">
          <NumberInput value={form.gps_north ?? ""} onChange={(e) => update("gps_north", num(e.target.value))} />
        </Field>
        <Field am="ከፍታ (ሜትር)" en="Elevation (m)">
          <NumberInput value={form.elevation_m ?? ""} onChange={(e) => update("elevation_m", num(e.target.value))} />
        </Field>
      </FormSection>

      {/* Section 3: Ownership & Construction */}
      <FormSection num={3} am="ባለቤትነት እና ግንባታ" en="Ownership & Construction" open={openSections[3]} onToggle={() => toggle(3)}>
        <Field am="የባለቤት ዓይነት" en="Owner Type" required>
          <Select
            value={form.owner_type ?? ""}
            onChange={(e) => update("owner_type", e.target.value || undefined)}
            options={OWNER_TYPES}
            placeholderAm="—"
          />
        </Field>
        <Field am="የባለቤት ስም" en="Owner Name" span={2}>
          <TextInput value={form.owner_name ?? ""} onChange={(e) => update("owner_name", e.target.value)} />
        </Field>
        <Field am="የገነባው" en="Built By">
          <TextInput value={form.built_by ?? ""} onChange={(e) => update("built_by", e.target.value)} />
        </Field>
        <Field am="የግንባታ ጊዜ" en="Construction Period" required>
          <TextInput value={form.construction_period ?? ""} onChange={(e) => update("construction_period", e.target.value)} />
        </Field>
        <Field am="የዕድሜ ዘዴ" en="Age Method" required>
          <Select
            value={form.age_method ?? ""}
            onChange={(e) => update("age_method", e.target.value || undefined)}
            options={AGE_METHODS}
            placeholderAm="—"
          />
        </Field>
      </FormSection>

      {/* Section 4: Dimensions & Building Details */}
      <FormSection num={4} am="ልኬቶችና የግንባታ ዝርዝር" en="Dimensions & Building Details" open={openSections[4]} onToggle={() => toggle(4)}>
        <Field am="ቁመት (ሜትር)" en="Height (m)">
          <NumberInput value={form.height_m ?? ""} onChange={(e) => update("height_m", num(e.target.value))} />
        </Field>
        <Field am="ርዝመት (ሜትር)" en="Length (m)">
          <NumberInput value={form.length_m ?? ""} onChange={(e) => update("length_m", num(e.target.value))} />
        </Field>
        <Field am="ስፋት (ሜትር)" en="Width (m)">
          <NumberInput value={form.width_m ?? ""} onChange={(e) => update("width_m", num(e.target.value))} />
        </Field>
        <Field am="የበሮች ብዛት" en="Doors">
          <NumberInput value={form.num_doors ?? ""} onChange={(e) => update("num_doors", num(e.target.value))} />
        </Field>
        <Field am="የመስኮቶች ብዛት" en="Windows">
          <NumberInput value={form.num_windows ?? ""} onChange={(e) => update("num_windows", num(e.target.value))} />
        </Field>
        <Field am="የክፍሎች ብዛት" en="Rooms">
          <NumberInput value={form.num_rooms ?? ""} onChange={(e) => update("num_rooms", num(e.target.value))} />
        </Field>
        <Field am="ቁሳቁስ" en="Material" span={3}>
          <TextInput value={form.material ?? ""} onChange={(e) => update("material", e.target.value)} />
        </Field>
        <Field am="የሐረሪ ቤት ደረጃ" en="Harari House Grades" span={2}>
          <CheckboxList
            options={HARARI_HOUSE_GRADES}
            value={form.harari_house_grades ?? []}
            onChange={(v) => update("harari_house_grades", v)}
          />
        </Field>
        <Field am="የመንደር ዓይነት" en="Neighborhood">
          <Select
            value={form.neighborhood_type ?? ""}
            onChange={(e) => update("neighborhood_type", e.target.value || undefined)}
            options={NEIGHBORHOOD_TYPES}
            placeholderAm="—"
          />
        </Field>
        <Field am="መግለጫ" en="Description" span={3}>
          <TextArea value={form.description ?? ""} onChange={(e) => update("description", e.target.value)} />
        </Field>
      </FormSection>

      {/* Section 5: Condition & Damage */}
      <FormSection num={5} am="ሁኔታ እና ጉዳት" en="Condition & Damage" open={openSections[5]} onToggle={() => toggle(5)}>
        <Field am="አጠቃላይ ሁኔታ" en="Overall Condition" span={3}>
          <Select
            value={form.overall_condition ?? ""}
            onChange={(e) => update("overall_condition", e.target.value || undefined)}
            options={OVERALL_CONDITIONS}
            placeholderAm="—"
          />
        </Field>
        {[
          { k: "damage_roof", am: "ጣሪያ", en: "Roof" },
          { k: "damage_cornice", am: "ኮርኒስ", en: "Cornice" },
          { k: "damage_wall", am: "ግድግዳ", en: "Wall" },
          { k: "damage_floor", am: "ወለል", en: "Floor" },
          { k: "damage_door", am: "በር", en: "Door" },
          { k: "damage_cupboard", am: "ካቦርድ", en: "Cupboard" },
          { k: "damage_upper_floor", am: "የላይኛው ወለል", en: "Upper Floor" },
          { k: "damage_dera", am: "ደራ", en: "Dera" },
          { k: "damage_pillar", am: "ምሰሶ", en: "Pillar" },
        ].map(({ k, am, en }) => (
          <Field key={k} am={am} en={en}>
            <Select
              value={(form[k as keyof ImmovableRecordInput] as string) ?? ""}
              onChange={(e) =>
                update(k as keyof ImmovableRecordInput, (e.target.value || undefined) as never)
              }
              options={DAMAGE_LEVELS}
              placeholderAm="—"
            />
          </Field>
        ))}
      </FormSection>

      {/* Section 6: Heritage Values & Threats */}
      <FormSection num={6} am="የቅርስ እሴት እና ስጋት" en="Heritage Values & Threats" open={openSections[6]} onToggle={() => toggle(6)}>
        {[
          { k: "value_historical", am: "ታሪካዊ እሴት", en: "Historical Value" },
          { k: "value_craftsmanship", am: "የእጅ ሥራ", en: "Craftsmanship" },
          { k: "value_artistic", am: "ሥነ ጥበባዊ", en: "Artistic" },
          { k: "value_scientific", am: "ሳይንሳዊ", en: "Scientific" },
          { k: "value_cultural", am: "ባህላዊ", en: "Cultural" },
        ].map(({ k, am, en }) => (
          <Field key={k} am={am} en={en} span={3}>
            <TextArea
              rows={2}
              value={(form[k as keyof ImmovableRecordInput] as string) ?? ""}
              onChange={(e) =>
                update(k as keyof ImmovableRecordInput, e.target.value as never)
              }
            />
          </Field>
        ))}
        <Field am="ስጋት አለ?" en="Has Threat?">
          <Switch
            checked={!!form.has_threat}
            onChange={(v) => update("has_threat", v)}
            am="አዎ"
            en="Yes"
          />
        </Field>
        {form.has_threat && (
          <Field am="የስጋት ምክንያት" en="Threat Reason" span={2}>
            <TextInput
              value={form.maintenance_reason ?? ""}
              onChange={(e) => update("maintenance_reason", e.target.value)}
            />
          </Field>
        )}
        <Field am="ጥገና ተደርጓል?" en="Maintenance Done?">
          <Switch
            checked={!!form.maintenance_done}
            onChange={(v) => update("maintenance_done", v)}
            am="አዎ"
            en="Yes"
          />
        </Field>
        <Field am="ጥገናው የተደረገበት" en="Maintenance By">
          <TextInput value={form.maintenance_by ?? ""} onChange={(e) => update("maintenance_by", e.target.value)} />
        </Field>
        <Field am="የጥገና ቀን" en="Maintenance Date">
          <TextInput
            type="date"
            value={form.maintenance_date ?? ""}
            onChange={(e) => update("maintenance_date", e.target.value)}
          />
        </Field>
        <Field am="የጥገና ብዛት" en="Maintenance Count">
          <NumberInput value={form.maintenance_count ?? ""} onChange={(e) => update("maintenance_count", num(e.target.value))} />
        </Field>
      </FormSection>

      {/* Section 7: Accessibility & Documentation */}
      <FormSection num={7} am="ተደራሽነት እና ሰነድ" en="Accessibility & Documentation" open={openSections[7]} onToggle={() => toggle(7)}>
        <Field am="የመከላከያ ደረጃ" en="Preventive Level">
          <Select
            value={form.preventive_level ?? ""}
            onChange={(e) => update("preventive_level", e.target.value || undefined)}
            options={QUALITY_LEVELS}
            placeholderAm="—"
          />
        </Field>
        <Field am="ተደራሽነት" en="Accessibility" span={2}>
          <Select
            value={form.accessibility ?? ""}
            onChange={(e) => update("accessibility", e.target.value || undefined)}
            options={ACCESSIBILITY_LEVELS}
            placeholderAm="—"
          />
        </Field>
        <Field am="ማስታወሻ" en="Notes" span={3}>
          <TextArea value={form.notes ?? ""} onChange={(e) => update("notes", e.target.value)} />
        </Field>
        <Field am="ተዛማጅ መረጃዎች" en="Related Documents" span={3}>
          <CheckboxList
            options={RELATED_DOCS}
            value={form.related_docs ?? []}
            onChange={(v) => update("related_docs", v)}
          />
        </Field>
      </FormSection>

      {/* Section 8: Oral History & Informant */}
      <FormSection num={8} am="የቃል ታሪክ እና መረጃ ሰጪ" en="Oral History & Informant" open={openSections[8]} onToggle={() => toggle(8)}>
        <Field am="የቃል ታሪክ አለ?" en="Has Oral History?">
          <Switch
            checked={!!form.has_oral_history}
            onChange={(v) => update("has_oral_history", v)}
            am="አዎ"
            en="Yes"
          />
        </Field>
        <Field am="የጠባቂ ስም" en="Caretaker Name">
          <TextInput value={form.caretaker_name ?? ""} onChange={(e) => update("caretaker_name", e.target.value)} />
        </Field>
        <Field am="የጠባቂ ሚና" en="Caretaker Role">
          <TextInput value={form.caretaker_role ?? ""} onChange={(e) => update("caretaker_role", e.target.value)} />
        </Field>
        <Field am="የመረጃ ሰጪ ስም" en="Informant Name">
          <TextInput value={form.informant_name ?? ""} onChange={(e) => update("informant_name", e.target.value)} />
        </Field>
        <Field am="ጾታ" en="Sex">
          <Select
            value={form.informant_sex ?? ""}
            onChange={(e) => update("informant_sex", e.target.value || undefined)}
            options={SEX_TYPES}
            placeholderAm="—"
          />
        </Field>
        <Field am="ዕድሜ" en="Age">
          <NumberInput value={form.informant_age ?? ""} onChange={(e) => update("informant_age", num(e.target.value))} />
        </Field>
        <Field am="የመዝጋቢ ቀን" en="Registrar Date">
          <TextInput
            type="date"
            value={form.registrar_date ?? ""}
            onChange={(e) => update("registrar_date", e.target.value)}
          />
        </Field>
      </FormSection>

      {/* Photos — only after record exists */}
      {mode === "edit" && initialRecord && (
        <PhotoUploader
          recordId={initialRecord.id}
          photos={photos}
          disabled={initialRecord.status !== "draft" && initialRecord.status !== "returned"}
          onUpload={(file) => uploadImmovablePhoto(initialRecord.id, file)}
          onDelete={(photoId) => deleteImmovablePhoto(initialRecord.id, photoId)}
          onPhotosChange={() => {
            qc.invalidateQueries({ queryKey: ["immovable", initialRecord.id] });
            qc.invalidateQueries({
              queryKey: ["record-detail", "immovable", initialRecord.id],
            });
          }}
        />
      )}

      {/* Error */}
      {error && (
        <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </div>
      )}

      {/* Sticky action bar */}
      <div className="sticky bottom-0 -mx-4 border-t border-border bg-background/95 px-4 py-3 backdrop-blur supports-[backdrop-filter]:bg-background/80">
        <div className="flex flex-wrap items-center justify-end gap-2">
          <button
            type="button"
            onClick={() => navigate({ to: "/records" })}
            className="font-amharic rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground hover:bg-accent"
          >
            {t("common.cancel")}
          </button>
          <button
            type="submit"
            disabled={saveMut.isPending}
            className="font-amharic inline-flex items-center gap-1.5 rounded-md border border-input bg-background px-3 py-2 text-sm font-medium text-foreground hover:bg-accent disabled:opacity-60"
          >
            {saveMut.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
            {t("forms.saveDraft")}
          </button>
          {mode === "edit" && initialRecord?.status !== "approved" && (
            <button
              type="button"
              onClick={handleSubmit}
              disabled={submitMut.isPending || saveMut.isPending}
              className="font-amharic inline-flex items-center gap-1.5 rounded-md bg-primary px-3 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 disabled:opacity-60"
            >
              {submitMut.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Send className="h-4 w-4" />}
              {t("forms.submitForReview")}
            </button>
          )}
        </div>
      </div>
    </form>
  );
}
