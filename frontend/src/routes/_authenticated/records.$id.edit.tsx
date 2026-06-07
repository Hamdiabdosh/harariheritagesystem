import { createFileRoute, Link } from "@tanstack/react-router";
import { useQuery } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { ArrowLeft } from "lucide-react";
import { getImmovable } from "@/api/immovable";
import { getMovable } from "@/api/movable";
import { ImmovableForm } from "@/components/forms/ImmovableForm";
import { MovableForm } from "@/components/forms/MovableForm";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { StatusBadge } from "@/components/common/StatusBadge";
import { StatusTimeline } from "@/components/common/StatusTimeline";
import { CommentThread } from "@/components/common/CommentThread";
import { PhotoGrid } from "@/components/common/PhotoGrid";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useAuthStore } from "@/stores/authStore";
import type {
  RecordType,
  StatusHistoryEntry,
  RecordComment,
  RecordPhoto,
} from "@/types";

export const Route = createFileRoute("/_authenticated/records/$id/edit")({
  component: EditRecordPage,
});

function EditRecordPage() {
  const { id } = Route.useParams();
  const { t } = useTranslation();
  const role = useAuthStore((s) => s.user?.role);

  const immovableQ = useQuery({
    queryKey: ["immovable", id],
    queryFn: () => getImmovable(id),
    retry: false,
  });

  const movableQ = useQuery({
    queryKey: ["movable", id],
    queryFn: () => getMovable(id),
    enabled: immovableQ.isError,
    retry: false,
  });

  const isLoading = immovableQ.isLoading || (immovableQ.isError && movableQ.isLoading);
  const isError = immovableQ.isError && movableQ.isError;
  const error = movableQ.error ?? immovableQ.error;

  const recordType: RecordType | null = immovableQ.data
    ? "immovable"
    : movableQ.data
      ? "movable"
      : null;

  const detail = immovableQ.data ?? movableQ.data;
  const canComment = role === "supervisor" || role === "manager";
  const canViewComments = true;

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center gap-3">
        <Link
          to="/records"
          className="font-amharic inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground"
        >
          <ArrowLeft className="h-4 w-4" />
          {t("detail.backToRecords")}
        </Link>
        {detail && (
          <>
            <span className="rounded-md bg-muted px-2 py-0.5 text-xs font-medium tabular-nums text-muted-foreground">
              {detail.record.record_id}
            </span>
            <StatusBadge status={detail.record.status} />
          </>
        )}
      </div>

      {isLoading && (
        <div className="flex justify-center rounded-xl border border-border bg-card p-12">
          <LoadingSpinner />
        </div>
      )}

      {isError && (
        <div className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-900">
          {(error as Error)?.message ?? t("records.loadError")}
        </div>
      )}

      {detail && recordType === "immovable" && (
        <>
          <ImmovableForm
            mode="edit"
            initialRecord={detail.record}
            photos={detail.photos}
          />
          <RecordDetailTabs
            recordType={recordType}
            recordId={id}
            history={detail.history}
            comments={detail.comments}
            photos={detail.photos}
            canComment={canComment}
            canViewComments={canViewComments}
            readOnly={role === "registrar"}
          />
        </>
      )}

      {detail && recordType === "movable" && (
        <>
          <MovableForm
            mode="edit"
            initialRecord={detail.record}
            photos={detail.photos}
          />
          <RecordDetailTabs
            recordType={recordType}
            recordId={id}
            history={detail.history}
            comments={detail.comments}
            photos={detail.photos}
            canComment={canComment}
            canViewComments={canViewComments}
            readOnly={role === "registrar"}
          />
        </>
      )}
    </div>
  );
}

function RecordDetailTabs({
  recordType,
  recordId,
  history,
  comments,
  photos,
  canComment,
  canViewComments,
  readOnly,
}: {
  recordType: RecordType;
  recordId: string;
  history: StatusHistoryEntry[];
  comments: RecordComment[];
  photos: RecordPhoto[];
  canComment: boolean;
  canViewComments: boolean;
  readOnly?: boolean;
}) {
  const { t } = useTranslation();

  return (
    <div className="rounded-xl border border-border bg-card p-4">
      <Tabs defaultValue="history">
        <TabsList>
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
        <TabsContent value="history" className="mt-4">
          <StatusTimeline history={history} />
        </TabsContent>
        <TabsContent value="comments" className="mt-4">
          <CommentThread
            comments={comments}
            recordType={recordType}
            recordId={recordId}
            canComment={canComment}
            canViewComments={canViewComments}
            readOnly={readOnly}
          />
        </TabsContent>
        <TabsContent value="photos" className="mt-4">
          <PhotoGrid photos={photos} />
        </TabsContent>
      </Tabs>
    </div>
  );
}
