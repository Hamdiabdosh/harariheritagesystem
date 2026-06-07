import type { ReactNode } from "react";
import { useTranslation } from "react-i18next";
import type { ImmovableRecord, MovableRecord, RecordComment, RecordPhoto, RecordType, StatusHistoryEntry } from "@/types";
import { StatusTimeline } from "@/components/common/StatusTimeline";
import { CommentThread } from "@/components/common/CommentThread";
import { PhotoGrid } from "@/components/common/PhotoGrid";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

export function ReadField({
  labelAm,
  labelEn,
  value,
}: {
  labelAm: string;
  labelEn: string;
  value: ReactNode;
}) {
  if (
    value === null ||
    value === undefined ||
    value === "" ||
    (Array.isArray(value) && value.length === 0)
  ) {
    return null;
  }
  return (
    <div>
      <div className="text-xs text-muted-foreground">
        <span className="font-amharic">{labelAm}</span> / {labelEn}
      </div>
      <div className="font-amharic mt-0.5 text-sm text-foreground">
        {Array.isArray(value) ? value.join(", ") : value}
      </div>
    </div>
  );
}

export type FieldDef = { key: string; am: string; en: string };

export const IMMOVABLE_FIELDS: FieldDef[] = [
  { key: "name_amharic", am: "ስም (አማርኛ)", en: "Name (Amharic)" },
  { key: "name_local", am: "ሌላ ስም", en: "Local Name" },
  { key: "category", am: "ምድብ", en: "Category" },
  { key: "current_use", am: "የአሁን አጠቃቀም", en: "Current Use" },
  { key: "current_use_other", am: "ሌላ አጠቃቀም", en: "Other Use" },
  { key: "previous_id", am: "የቀድሞ መለያ", en: "Previous ID" },
  { key: "woreda", am: "ወረዳ", en: "Woreda" },
  { key: "kebele", am: "ቀበሌ", en: "Kebele" },
  { key: "house_number", am: "የቤት ቁጥር", en: "House Number" },
  { key: "street_number", am: "የመንገድ ቁጥር", en: "Street Number" },
  { key: "gate", am: "በር", en: "Gate" },
  { key: "gps_east", am: "GPS ምስራቅ", en: "GPS East" },
  { key: "gps_north", am: "GPS ሰሜን", en: "GPS North" },
  { key: "elevation_m", am: "ከፍታ (ሜ)", en: "Elevation (m)" },
  { key: "owner_type", am: "የባለቤት ዓይነት", en: "Owner Type" },
  { key: "owner_name", am: "የባለቤት ስም", en: "Owner Name" },
  { key: "map_reference", am: "የካርታ ማጣቀሻ", en: "Map Reference" },
  { key: "built_by", am: "የገነባው", en: "Built By" },
  { key: "construction_period", am: "የግንባታ ጊዜ", en: "Construction Period" },
  { key: "age_method", am: "የዕድሜ ዘዴ", en: "Age Method" },
  { key: "height_m", am: "ቁመት (ሜ)", en: "Height (m)" },
  { key: "length_m", am: "ርዝመት (ሜ)", en: "Length (m)" },
  { key: "width_m", am: "ስፋት (ሜ)", en: "Width (m)" },
  { key: "num_doors", am: "የበሮች ብዛት", en: "Doors" },
  { key: "num_windows", am: "የመስኮቶች ብዛት", en: "Windows" },
  { key: "num_rooms", am: "የክፍሎች ብዛት", en: "Rooms" },
  { key: "material", am: "ቁሳቁስ", en: "Material" },
  { key: "harari_house_grades", am: "የሐረሪ ቤት ደረጃ", en: "Harari House Grades" },
  { key: "neighborhood_type", am: "የመንደር ዓይነት", en: "Neighborhood" },
  { key: "description", am: "መግለጫ", en: "Description" },
  { key: "overall_condition", am: "አጠቃላይ ሁኔታ", en: "Overall Condition" },
  { key: "damage_roof", am: "ጣሪያ", en: "Roof Damage" },
  { key: "damage_cornice", am: "ኮርኒስ", en: "Cornice Damage" },
  { key: "damage_wall", am: "ግድግዳ", en: "Wall Damage" },
  { key: "damage_floor", am: "ወለል", en: "Floor Damage" },
  { key: "damage_door", am: "በር", en: "Door Damage" },
  { key: "damage_cupboard", am: "ካቦርድ", en: "Cupboard Damage" },
  { key: "damage_upper_floor", am: "የላይኛው ወለል", en: "Upper Floor Damage" },
  { key: "damage_dera", am: "ደራ", en: "Dera Damage" },
  { key: "damage_pillar", am: "ምሰሶ", en: "Pillar Damage" },
  { key: "value_historical", am: "ታሪካዊ እሴት", en: "Historical Value" },
  { key: "value_craftsmanship", am: "የእጅ ሥራ", en: "Craftsmanship" },
  { key: "value_artistic", am: "ሥነ ጥበባዊ", en: "Artistic Value" },
  { key: "value_scientific", am: "ሳይንሳዊ", en: "Scientific Value" },
  { key: "value_cultural", am: "ባህላዊ", en: "Cultural Value" },
  { key: "has_threat", am: "ስጋት", en: "Has Threat" },
  { key: "maintenance_done", am: "ጥገና", en: "Maintenance Done" },
  { key: "maintenance_reason", am: "የጥገና ምክንያት", en: "Maintenance Reason" },
  { key: "maintenance_by", am: "ጥገናው የተደረገበት", en: "Maintenance By" },
  { key: "maintenance_count", am: "የጥገና ብዛት", en: "Maintenance Count" },
  { key: "preventive_level", am: "የመከላከያ ደረጃ", en: "Preventive Level" },
  { key: "accessibility", am: "ተደራሽነት", en: "Accessibility" },
  { key: "notes", am: "ማስታወሻ", en: "Notes" },
  { key: "caretaker_name", am: "የጠባቂ ስም", en: "Caretaker Name" },
  { key: "caretaker_role", am: "የጠባቂ ሚና", en: "Caretaker Role" },
  { key: "informant_name", am: "የመረጃ ሰጪ", en: "Informant Name" },
  { key: "informant_sex", am: "ጾታ", en: "Sex" },
  { key: "informant_age", am: "ዕድሜ", en: "Age" },
  { key: "registrar_date", am: "የመዝጋቢ ቀን", en: "Registrar Date" },
];

export const MOVABLE_FIELDS: FieldDef[] = [
  { key: "name_amharic", am: "ስም (አማርኛ)", en: "Name (Amharic)" },
  { key: "name_local", am: "ሌላ ስም", en: "Local Name" },
  { key: "category", am: "ምድብ", en: "Category" },
  { key: "current_use", am: "የአሁን አጠቃቀም", en: "Current Use" },
  { key: "previous_id", am: "የቀድሞ መለያ", en: "Previous ID" },
  { key: "location_name", am: "የቦታ ስም", en: "Location Name" },
  { key: "woreda", am: "ወረዳ", en: "Woreda" },
  { key: "kebele", am: "ቀበሌ", en: "Kebele" },
  { key: "house_number", am: "የቤት ቁጥር", en: "House Number" },
  { key: "owner_type", am: "የባለቤት ዓይነት", en: "Owner Type" },
  { key: "owner_name", am: "የባለቤት ስም", en: "Owner Name" },
  { key: "storage_location", am: "ማስቀመጫ", en: "Storage Location" },
  { key: "storage_location_other", am: "ሌላ ማስቀመጫ", en: "Other Storage" },
  { key: "made_by", am: "የተሰራው", en: "Made By" },
  { key: "period_made", am: "የተሰራበት ጊዜ", en: "Period Made" },
  { key: "age_method", am: "የዕድሜ ዘዴ", en: "Age Method" },
  { key: "acquisition_methods", am: "የግዢ ዘዴ", en: "Acquisition Methods" },
  { key: "height_cm", am: "ቁመት (ሴ.ሜ)", en: "Height (cm)" },
  { key: "width_cm", am: "ስፋት (ሴ.ሜ)", en: "Width (cm)" },
  { key: "length_cm", am: "ርዝመት (ሴ.ሜ)", en: "Length (cm)" },
  { key: "diameter_cm", am: "ዲያሜትር", en: "Diameter (cm)" },
  { key: "thickness_cm", am: "ውፍረት", en: "Thickness (cm)" },
  { key: "weight_kg", am: "ክብደት (ኪ.ግ)", en: "Weight (kg)" },
  { key: "num_pages", am: "ገጾች", en: "Pages" },
  { key: "num_chapters", am: "ምዕራፎች", en: "Chapters" },
  { key: "num_illustrations", am: "ምስሎች", en: "Illustrations" },
  { key: "color_type", am: "ቀለም", en: "Color" },
  { key: "has_decoration", am: "ጌጣጌጥ", en: "Has Decoration" },
  { key: "materials", am: "ቁሳቁሶች", en: "Materials" },
  { key: "material_other", am: "ሌላ ቁሳቁስ", en: "Other Material" },
  { key: "description", am: "መግለጫ", en: "Description" },
  { key: "notable_because", am: "ታዋቂ በ", en: "Notable Because" },
  { key: "notable_other", am: "ሌላ", en: "Other Notable" },
  { key: "significance", am: "ጠቀሜታ", en: "Significance" },
  { key: "condition", am: "ሁኔታ", en: "Condition" },
  { key: "has_threat", am: "ስጋት", en: "Has Threat" },
  { key: "threat_description", am: "የስጋት መግለጫ", en: "Threat Description" },
  { key: "maintenance_done", am: "ጥገና", en: "Maintenance Done" },
  { key: "maintenance_by", am: "ጥገናው የተደረገበት", en: "Maintenance By" },
  { key: "maintenance_date", am: "የጥገና ቀን", en: "Maintenance Date" },
  { key: "maintenance_count", am: "የጥገና ብዛት", en: "Maintenance Count" },
  { key: "preventive_level", am: "የመከላከያ ደረጃ", en: "Preventive Level" },
  { key: "accessibility", am: "ተደራሽነት", en: "Accessibility" },
  { key: "notes", am: "ማስታወሻ", en: "Notes" },
  { key: "informant_name", am: "የመረጃ ሰጪ", en: "Informant Name" },
  { key: "informant_sex", am: "ጾታ", en: "Sex" },
  { key: "informant_age", am: "ዕድሜ", en: "Age" },
  { key: "informant_occupation", am: "ሙያ", en: "Occupation" },
  { key: "caretaker_name", am: "የጠባቂ ስም", en: "Caretaker Name" },
  { key: "caretaker_role", am: "የጠባቂ ሚና", en: "Caretaker Role" },
  { key: "registrar_date", am: "የመዝጋቢ ቀን", en: "Registrar Date" },
];

function formatValue(value: unknown, t: (k: string) => string): ReactNode {
  if (typeof value === "boolean") return value ? t("common.yes") : t("common.no");
  return value as ReactNode;
}

export function RecordDetailsGrid({
  record,
  fields,
}: {
  record: ImmovableRecord | MovableRecord;
  fields: FieldDef[];
}) {
  const { t } = useTranslation();
  const rec = record as Record<string, unknown>;
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
      {fields.map(({ key, am, en }) => (
        <ReadField
          key={key}
          labelAm={am}
          labelEn={en}
          value={formatValue(rec[key], t)}
        />
      ))}
    </div>
  );
}

export function RecordDetailTabs({
  recordType,
  recordId,
  record,
  history,
  comments,
  photos,
  canComment,
}: {
  recordType: RecordType;
  recordId: string;
  record: ImmovableRecord | MovableRecord;
  history: StatusHistoryEntry[];
  comments: RecordComment[];
  photos: RecordPhoto[];
  canComment: boolean;
}) {
  const { t } = useTranslation();
  const fields = recordType === "immovable" ? IMMOVABLE_FIELDS : MOVABLE_FIELDS;

  return (
    <div className="rounded-xl border border-border bg-card p-4">
      <Tabs defaultValue="details">
        <TabsList>
          <TabsTrigger value="details" className="font-amharic">
            {t("detail.tabs.details")}
          </TabsTrigger>
          <TabsTrigger value="history" className="font-amharic">
            {t("detail.tabs.history")}
          </TabsTrigger>
          <TabsTrigger value="comments" className="font-amharic">
            {t("detail.tabs.comments")}
          </TabsTrigger>
          <TabsTrigger value="photos" className="font-amharic">
            {t("detail.tabs.photos")}
          </TabsTrigger>
        </TabsList>
        <TabsContent value="details" className="mt-4">
          <RecordDetailsGrid record={record} fields={fields} />
        </TabsContent>
        <TabsContent value="history" className="mt-4">
          <StatusTimeline history={history} />
        </TabsContent>
        <TabsContent value="comments" className="mt-4">
          <CommentThread
            comments={comments}
            recordType={recordType}
            recordId={recordId}
            canComment={canComment}
          />
        </TabsContent>
        <TabsContent value="photos" className="mt-4">
          <PhotoGrid photos={photos} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
