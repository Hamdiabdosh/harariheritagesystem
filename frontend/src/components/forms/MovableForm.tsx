import { useState, type FormEvent } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { Save, Send, Loader2 } from "lucide-react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { MovableRecord, MovableRecordInput, RecordPhoto } from "@/types";
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
  createMovable,
  updateMovable,
  submitMovable,
  uploadMovablePhoto,
  deleteMovablePhoto,
} from "@/api/movable";
import { toDateInputValue } from "@/lib/recordPayload";
import { validateMovableForSubmit } from "@/lib/recordSubmitValidation";
import { getApiErrorMessage } from "@/lib/apiError";
import {
  MOVABLE_CATEGORIES,
  MOVABLE_OWNER_TYPES,
  STORAGE_LOCATIONS,
  MOVABLE_AGE_METHODS,
  ACQUISITION_METHODS,
  MOVABLE_CONDITIONS,
  MATERIALS,
  NOTABLE_BECAUSE,
  QUALITY_LEVELS,
  ACCESSIBILITY_LEVELS,
  SEX_TYPES,
} from "./movableOptions";

interface MovableFormProps {
  mode: "create" | "edit";
  initialRecord?: MovableRecord;
  photos?: RecordPhoto[];
}

const num = (v: string) => (v === "" ? undefined : Number(v));

export function MovableForm({ mode, initialRecord, photos = [] }: MovableFormProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const qc = useQueryClient();

  const [form, setForm] = useState<MovableRecordInput>(() => ({
    name_amharic: initialRecord?.name_amharic ?? "",
    name_local: initialRecord?.name_local ?? "",
    category: initialRecord?.category ?? "",
    current_use: initialRecord?.current_use ?? "",
    previous_id: initialRecord?.previous_id ?? "",
    location_name: initialRecord?.location_name ?? "",
    woreda: initialRecord?.woreda ?? "",
    kebele: initialRecord?.kebele ?? "",
    house_number: initialRecord?.house_number ?? "",
    owner_type: initialRecord?.owner_type,
    owner_name: initialRecord?.owner_name ?? "",
    storage_location: initialRecord?.storage_location,
    storage_location_other: initialRecord?.storage_location_other ?? "",
    made_by: initialRecord?.made_by ?? "",
    period_made: initialRecord?.period_made ?? "",
    age_method: initialRecord?.age_method,
    acquisition_methods: initialRecord?.acquisition_methods ?? [],
    height_cm: initialRecord?.height_cm,
    width_cm: initialRecord?.width_cm,
    length_cm: initialRecord?.length_cm,
    diameter_cm: initialRecord?.diameter_cm,
    thickness_cm: initialRecord?.thickness_cm,
    weight_kg: initialRecord?.weight_kg,
    num_pages: initialRecord?.num_pages,
    num_chapters: initialRecord?.num_chapters,
    num_illustrations: initialRecord?.num_illustrations,
    color_type: initialRecord?.color_type ?? "",
    has_decoration: initialRecord?.has_decoration ?? false,
    materials: initialRecord?.materials ?? [],
    material_other: initialRecord?.material_other ?? "",
    description: initialRecord?.description ?? "",
    notable_because: initialRecord?.notable_because ?? [],
    notable_other: initialRecord?.notable_other ?? "",
    significance: initialRecord?.significance ?? "",
    condition: initialRecord?.condition,
    has_threat: initialRecord?.has_threat ?? false,
    threat_description: initialRecord?.threat_description ?? "",
    maintenance_done: initialRecord?.maintenance_done ?? false,
    maintenance_by: initialRecord?.maintenance_by ?? "",
    maintenance_date: toDateInputValue(initialRecord?.maintenance_date),
    maintenance_count: initialRecord?.maintenance_count,
    preventive_level: initialRecord?.preventive_level,
    accessibility: initialRecord?.accessibility,
    notes: initialRecord?.notes ?? "",
    informant_name: initialRecord?.informant_name ?? "",
    informant_sex: initialRecord?.informant_sex,
    informant_age: initialRecord?.informant_age,
    informant_occupation: initialRecord?.informant_occupation ?? "",
    caretaker_name: initialRecord?.caretaker_name ?? "",
    caretaker_role: initialRecord?.caretaker_role ?? "",
    registrar_date: toDateInputValue(initialRecord?.registrar_date),
  }));

  const [openSections, setOpenSections] = useState<Record<number, boolean>>({
    1: true,
    2: true,
    3: false,
    4: false,
    5: false,
    6: false,
    7: false,
  });
  const toggle = (n: number) => setOpenSections((s) => ({ ...s, [n]: !s[n] }));

  const [error, setError] = useState<string | null>(null);

  const update = <K extends keyof MovableRecordInput>(
    key: K,
    value: MovableRecordInput[K],
  ) => setForm((f) => ({ ...f, [key]: value }));

  const saveMut = useMutation({
    mutationFn: async (body: MovableRecordInput) => {
      if (mode === "edit" && initialRecord) {
        await updateMovable(initialRecord.id, body);
        return initialRecord.id;
      }
      const result = await createMovable(body);
      return result.id;
    },
    onSuccess: (id) => {
      qc.invalidateQueries({ queryKey: ["records"] });
      qc.invalidateQueries({ queryKey: ["movable", id] });
      qc.invalidateQueries({ queryKey: ["record-detail", "movable", id] });
    },
    onError: (e: unknown) => setError(getApiErrorMessage(e)),
  });

  const submitMut = useMutation({
    mutationFn: async () => {
      const id =
        mode === "edit" && initialRecord
          ? initialRecord.id
          : await saveMut.mutateAsync(form);
      if (mode === "edit" && initialRecord) {
        await updateMovable(initialRecord.id, form);
      }
      await submitMovable(id);
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
    const id = await saveMut.mutateAsync(form).catch(() => null);
    if (!id) return;
    if (mode === "create") navigate({ to: "/records/$id/edit", params: { id } });
  };

  const handleSubmit = () => {
    setError(null);
    const missing = validateMovableForSubmit(form);
    if (missing.length > 0) {
      setOpenSections((s) => ({
        ...s,
        1: s[1] || missing.some((f) => ["name_amharic", "category"].includes(f)),
        2:
          s[2] ||
          missing.some((f) =>
            ["woreda", "kebele", "owner_type", "storage_location"].includes(f),
          ),
        4: s[4] || missing.includes("materials"),
      }));
      const labels = missing.map((f) => t(`forms.fields.${f}`)).join(", ");
      setError(t("forms.submitMissing", { fields: labels }));
      return;
    }
    submitMut.mutate();
  };

  const canEditPhotos =
    initialRecord &&
    (initialRecord.status === "draft" || initialRecord.status === "returned");

  return (
    <form onSubmit={handleSave} className="space-y-4 pb-24">
      <div className="rounded-xl border border-primary/20 bg-primary/5 p-4">
        <h2 className="font-amharic text-lg font-bold text-foreground">
          {t("movable.title")}
        </h2>
        <p className="font-amharic mt-1 text-xs text-muted-foreground">
          {t("movable.subtitle")}
        </p>
        {initialRecord && (
          <p className="mt-2 text-xs text-muted-foreground tabular-nums">
            ID: {initialRecord.record_id} · {t(`status.${initialRecord.status}`)}
          </p>
        )}
      </div>

      <FormSection
        num={1}
        am="መለያ"
        en="Identification"
        open={openSections[1]}
        onToggle={() => toggle(1)}
      >
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
        <Field am="ምድብ" en="Category" required>
          <Select
            value={form.category ?? ""}
            onChange={(e) => update("category", e.target.value || undefined)}
            options={MOVABLE_CATEGORIES}
            placeholderAm="—"
          />
        </Field>
        <Field am="የአሁን አጠቃቀም" en="Current Use">
          <TextInput
            value={form.current_use ?? ""}
            onChange={(e) => update("current_use", e.target.value)}
          />
        </Field>
        <Field am="የቀድሞ መለያ" en="Previous ID">
          <TextInput
            value={form.previous_id ?? ""}
            onChange={(e) => update("previous_id", e.target.value)}
          />
        </Field>
      </FormSection>

      <FormSection
        num={2}
        am="አካባቢ እና ማስቀመጫ"
        en="Location & Storage"
        open={openSections[2]}
        onToggle={() => toggle(2)}
      >
        <Field am="የቦታ ስም" en="Location Name" span={2}>
          <TextInput
            value={form.location_name ?? ""}
            onChange={(e) => update("location_name", e.target.value)}
          />
        </Field>
        <Field am="ወረዳ" en="Woreda">
          <TextInput
            value={form.woreda ?? ""}
            onChange={(e) => update("woreda", e.target.value)}
          />
        </Field>
        <Field am="ቀበሌ" en="Kebele">
          <TextInput
            value={form.kebele ?? ""}
            onChange={(e) => update("kebele", e.target.value)}
          />
        </Field>
        <Field am="የቤት ቁጥር" en="House Number">
          <TextInput
            value={form.house_number ?? ""}
            onChange={(e) => update("house_number", e.target.value)}
          />
        </Field>
        <Field am="የባለቤት ዓይነት" en="Owner Type">
          <Select
            value={form.owner_type ?? ""}
            onChange={(e) => update("owner_type", e.target.value || undefined)}
            options={MOVABLE_OWNER_TYPES}
            placeholderAm="—"
          />
        </Field>
        <Field am="የባለቤት ስም" en="Owner Name" span={2}>
          <TextInput
            value={form.owner_name ?? ""}
            onChange={(e) => update("owner_name", e.target.value)}
          />
        </Field>
        <Field am="ማስቀመጫ" en="Storage Location">
          <Select
            value={form.storage_location ?? ""}
            onChange={(e) => update("storage_location", e.target.value || undefined)}
            options={STORAGE_LOCATIONS}
            placeholderAm="—"
          />
        </Field>
        <Field am="ሌላ ማስቀመጫ" en="Other Storage">
          <TextInput
            value={form.storage_location_other ?? ""}
            onChange={(e) => update("storage_location_other", e.target.value)}
          />
        </Field>
      </FormSection>

      <FormSection
        num={3}
        am="አመጣጥ እና ግዥ"
        en="Origin & Acquisition"
        open={openSections[3]}
        onToggle={() => toggle(3)}
      >
        <Field am="የተሰራው" en="Made By">
          <TextInput
            value={form.made_by ?? ""}
            onChange={(e) => update("made_by", e.target.value)}
          />
        </Field>
        <Field am="የተሰራበት ጊዜ" en="Period Made">
          <TextInput
            value={form.period_made ?? ""}
            onChange={(e) => update("period_made", e.target.value)}
          />
        </Field>
        <Field am="የዕድሜ ዘዴ" en="Age Method">
          <Select
            value={form.age_method ?? ""}
            onChange={(e) => update("age_method", e.target.value || undefined)}
            options={MOVABLE_AGE_METHODS}
            placeholderAm="—"
          />
        </Field>
        <Field am="የግዢ ዘዴ" en="Acquisition Methods" span={3}>
          <CheckboxList
            options={ACQUISITION_METHODS}
            value={form.acquisition_methods ?? []}
            onChange={(v) => update("acquisition_methods", v)}
          />
        </Field>
      </FormSection>

      <FormSection
        num={4}
        am="አካላዊ ገጽታ"
        en="Physical Description"
        open={openSections[4]}
        onToggle={() => toggle(4)}
      >
        <Field am="ቁመት (ሴ.ሜ)" en="Height (cm)">
          <NumberInput
            value={form.height_cm ?? ""}
            onChange={(e) => update("height_cm", num(e.target.value))}
          />
        </Field>
        <Field am="ስፋት (ሴ.ሜ)" en="Width (cm)">
          <NumberInput
            value={form.width_cm ?? ""}
            onChange={(e) => update("width_cm", num(e.target.value))}
          />
        </Field>
        <Field am="ርዝመት (ሴ.ሜ)" en="Length (cm)">
          <NumberInput
            value={form.length_cm ?? ""}
            onChange={(e) => update("length_cm", num(e.target.value))}
          />
        </Field>
        <Field am="ዲያሜትር (ሴ.ሜ)" en="Diameter (cm)">
          <NumberInput
            value={form.diameter_cm ?? ""}
            onChange={(e) => update("diameter_cm", num(e.target.value))}
          />
        </Field>
        <Field am="ውፍረት (ሴ.ሜ)" en="Thickness (cm)">
          <NumberInput
            value={form.thickness_cm ?? ""}
            onChange={(e) => update("thickness_cm", num(e.target.value))}
          />
        </Field>
        <Field am="ክብደት (ኪ.ግ)" en="Weight (kg)">
          <NumberInput
            value={form.weight_kg ?? ""}
            onChange={(e) => update("weight_kg", num(e.target.value))}
          />
        </Field>
        <Field am="ገጾች" en="Pages">
          <NumberInput
            value={form.num_pages ?? ""}
            onChange={(e) => update("num_pages", num(e.target.value))}
          />
        </Field>
        <Field am="ምዕራፎች" en="Chapters">
          <NumberInput
            value={form.num_chapters ?? ""}
            onChange={(e) => update("num_chapters", num(e.target.value))}
          />
        </Field>
        <Field am="ምስሎች" en="Illustrations">
          <NumberInput
            value={form.num_illustrations ?? ""}
            onChange={(e) => update("num_illustrations", num(e.target.value))}
          />
        </Field>
        <Field am="ቀለም" en="Color Type">
          <TextInput
            value={form.color_type ?? ""}
            onChange={(e) => update("color_type", e.target.value)}
          />
        </Field>
        <Field am="ጌጣጌጥ አለ?" en="Has Decoration?">
          <Switch
            checked={!!form.has_decoration}
            onChange={(v) => update("has_decoration", v)}
            am="አዎ"
            en="Yes"
          />
        </Field>
        <Field am="ቁሳቁሶች" en="Materials" span={3}>
          <CheckboxList
            options={MATERIALS}
            value={form.materials ?? []}
            onChange={(v) => update("materials", v)}
          />
        </Field>
        <Field am="ሌላ ቁሳቁስ" en="Other Material">
          <TextInput
            value={form.material_other ?? ""}
            onChange={(e) => update("material_other", e.target.value)}
          />
        </Field>
        <Field am="መግለጫ" en="Description" span={3}>
          <TextArea
            value={form.description ?? ""}
            onChange={(e) => update("description", e.target.value)}
          />
        </Field>
      </FormSection>

      <FormSection
        num={5}
        am="ጠቀሜታ"
        en="Significance"
        open={openSections[5]}
        onToggle={() => toggle(5)}
      >
        <Field am="ታዋቂ በ..." en="Notable Because" span={3}>
          <CheckboxList
            options={NOTABLE_BECAUSE}
            value={form.notable_because ?? []}
            onChange={(v) => update("notable_because", v)}
          />
        </Field>
        <Field am="ሌላ" en="Other Notable">
          <TextInput
            value={form.notable_other ?? ""}
            onChange={(e) => update("notable_other", e.target.value)}
          />
        </Field>
        <Field am="ጠቀሜታ" en="Significance" span={3}>
          <TextArea
            value={form.significance ?? ""}
            onChange={(e) => update("significance", e.target.value)}
          />
        </Field>
      </FormSection>

      <FormSection
        num={6}
        am="ሁኔታ እና ጥበቃ"
        en="Condition & Conservation"
        open={openSections[6]}
        onToggle={() => toggle(6)}
      >
        <Field am="ሁኔታ" en="Condition" span={3}>
          <Select
            value={form.condition ?? ""}
            onChange={(e) => update("condition", e.target.value || undefined)}
            options={MOVABLE_CONDITIONS}
            placeholderAm="—"
          />
        </Field>
        <Field am="ስጋት አለ?" en="Has Threat?">
          <Switch
            checked={!!form.has_threat}
            onChange={(v) => update("has_threat", v)}
            am="አዎ"
            en="Yes"
          />
        </Field>
        <Field am="የስጋት መግለጫ" en="Threat Description" span={2}>
          <TextInput
            value={form.threat_description ?? ""}
            onChange={(e) => update("threat_description", e.target.value)}
          />
        </Field>
        <Field am="ጥገና ተደርጓል?" en="Maintenance Done?">
          <Switch
            checked={!!form.maintenance_done}
            onChange={(v) => update("maintenance_done", v)}
            am="አዎ"
            en="Yes"
          />
        </Field>
        <Field am="ጥገናው የተደረገበት" en="Maintenance By">
          <TextInput
            value={form.maintenance_by ?? ""}
            onChange={(e) => update("maintenance_by", e.target.value)}
          />
        </Field>
        <Field am="የጥገና ቀን" en="Maintenance Date">
          <TextInput
            type="date"
            value={form.maintenance_date ?? ""}
            onChange={(e) => update("maintenance_date", e.target.value)}
          />
        </Field>
        <Field am="የጥገና ብዛት" en="Maintenance Count">
          <NumberInput
            value={form.maintenance_count ?? ""}
            onChange={(e) => update("maintenance_count", num(e.target.value))}
          />
        </Field>
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
          <TextArea
            value={form.notes ?? ""}
            onChange={(e) => update("notes", e.target.value)}
          />
        </Field>
      </FormSection>

      <FormSection
        num={7}
        am="መረጃ ሰጪ እና መዝጋቢ"
        en="Informant & Registrar"
        open={openSections[7]}
        onToggle={() => toggle(7)}
      >
        <Field am="የመረጃ ሰጪ ስም" en="Informant Name">
          <TextInput
            value={form.informant_name ?? ""}
            onChange={(e) => update("informant_name", e.target.value)}
          />
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
          <NumberInput
            value={form.informant_age ?? ""}
            onChange={(e) => update("informant_age", num(e.target.value))}
          />
        </Field>
        <Field am="ሙያ" en="Occupation">
          <TextInput
            value={form.informant_occupation ?? ""}
            onChange={(e) => update("informant_occupation", e.target.value)}
          />
        </Field>
        <Field am="የጠባቂ ስም" en="Caretaker Name">
          <TextInput
            value={form.caretaker_name ?? ""}
            onChange={(e) => update("caretaker_name", e.target.value)}
          />
        </Field>
        <Field am="የጠባቂ ሚና" en="Caretaker Role">
          <TextInput
            value={form.caretaker_role ?? ""}
            onChange={(e) => update("caretaker_role", e.target.value)}
          />
        </Field>
        <Field am="የመዝጋቢ ቀን" en="Registrar Date">
          <TextInput
            type="date"
            value={form.registrar_date ?? ""}
            onChange={(e) => update("registrar_date", e.target.value)}
          />
        </Field>
      </FormSection>

      {mode === "edit" && initialRecord && (
        <PhotoUploader
          recordId={initialRecord.id}
          photos={photos}
          disabled={!canEditPhotos}
          onUpload={(file) => uploadMovablePhoto(initialRecord.id, file)}
          onDelete={(photoId) => deleteMovablePhoto(initialRecord.id, photoId)}
          onPhotosChange={() => {
            qc.invalidateQueries({ queryKey: ["movable", initialRecord.id] });
            qc.invalidateQueries({
              queryKey: ["record-detail", "movable", initialRecord.id],
            });
          }}
        />
      )}

      {error && (
        <div className="rounded-xl border border-destructive/30 bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </div>
      )}

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
            {saveMut.isPending ? (
              <Loader2 className="h-4 w-4 animate-spin" />
            ) : (
              <Save className="h-4 w-4" />
            )}
            {t("forms.saveDraft")}
          </button>
          {mode === "edit" && initialRecord?.status !== "approved" && (
            <button
              type="button"
              onClick={handleSubmit}
              disabled={submitMut.isPending || saveMut.isPending}
              className="font-amharic inline-flex items-center gap-1.5 rounded-md bg-primary px-3 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 disabled:opacity-60"
            >
              {submitMut.isPending ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Send className="h-4 w-4" />
              )}
              {t("forms.submitForReview")}
            </button>
          )}
        </div>
      </div>
    </form>
  );
}
