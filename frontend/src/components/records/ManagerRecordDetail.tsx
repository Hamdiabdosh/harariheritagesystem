import { useState } from "react";
import { useNavigate } from "@tanstack/react-router";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { Loader2, CheckCircle2, RotateCcw, Download } from "lucide-react";
import { toast } from "sonner";
import { getImmovable } from "@/api/immovable";
import { getMovable } from "@/api/movable";
import { finalApprove, finalReturn } from "@/api/workflow";
import { exportPDF, downloadBlob } from "@/api/export";
import type { RecordType } from "@/types";
import { StatusBadge } from "@/components/common/StatusBadge";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { ReturnModal } from "@/components/common/ReturnModal";
import { RecordDetailTabs } from "./recordDetailShared";

interface ManagerRecordDetailProps {
  recordId: string;
  recordType: RecordType;
}

export function ManagerRecordDetail({ recordId, recordType }: ManagerRecordDetailProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const qc = useQueryClient();
  const [returnOpen, setReturnOpen] = useState(false);
  const [pdfLoading, setPdfLoading] = useState(false);

  const immovableQ = useQuery({
    queryKey: ["immovable", recordId],
    queryFn: () => getImmovable(recordId),
    enabled: recordType === "immovable",
  });

  const movableQ = useQuery({
    queryKey: ["movable", recordId],
    queryFn: () => getMovable(recordId),
    enabled: recordType === "movable",
  });

  const query = recordType === "immovable" ? immovableQ : movableQ;
  const detail = query.data;

  const invalidateAll = () => {
    qc.invalidateQueries({ queryKey: [recordType, recordId] });
    qc.invalidateQueries({ queryKey: ["dashboard", "stats"] });
    qc.invalidateQueries({ queryKey: ["records"] });
  };

  const approveMut = useMutation({
    mutationFn: () => finalApprove(recordType, recordId),
    onSuccess: () => {
      invalidateAll();
      toast.success(t("toast.finalApproveSuccess"));
      navigate({ to: "/reviewed" });
    },
    onError: () => toast.error(t("toast.error")),
  });

  const returnMut = useMutation({
    mutationFn: (comment: string) => finalReturn(recordType, recordId, comment),
    onSuccess: () => {
      invalidateAll();
      toast.success(t("toast.returnSuccess"));
      navigate({ to: "/reviewed" });
    },
    onError: () => toast.error(t("toast.error")),
  });

  const handleDownloadPdf = async () => {
    if (!detail) return;
    setPdfLoading(true);
    try {
      const blob = await exportPDF(recordType, recordId);
      downloadBlob(blob, `${detail.record.record_id}.pdf`);
    } catch {
      toast.error(t("toast.error"));
    } finally {
      setPdfLoading(false);
    }
  };

  if (query.isLoading) {
    return (
      <div className="flex justify-center rounded-xl border border-border bg-card p-12">
        <LoadingSpinner />
      </div>
    );
  }

  if (query.isError || !detail) {
    return (
      <div className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-900">
        {(query.error as Error)?.message ?? t("records.loadError")}
      </div>
    );
  }

  const { record, photos, comments, history } = detail;
  const busy = approveMut.isPending || returnMut.isPending;

  return (
    <div className="space-y-4">
      <div className="rounded-xl border border-border bg-card p-4">
        <div className="flex flex-wrap items-start justify-between gap-3">
          <div>
            <div className="flex flex-wrap items-center gap-2">
              <span className="rounded-md bg-muted px-2 py-0.5 text-xs font-medium tabular-nums">
                {record.record_id}
              </span>
              <StatusBadge status={record.status} />
            </div>
            <h2 className="font-amharic mt-2 text-xl font-bold text-foreground">
              {record.name_amharic}
            </h2>
            <p className="font-amharic mt-1 text-sm text-muted-foreground">
              {[record.woreda, record.kebele].filter(Boolean).join(" / ")}{" "}
              · {t(`recordType.${recordType}`)}
            </p>
          </div>
        </div>

        {(record.status === "under_review" || record.status !== "draft") && (
          <div className="mt-4 border-t border-border pt-4">
            {record.status === "under_review" && (
              <>
                <p className="font-amharic mb-3 text-sm font-semibold text-foreground">
                  {t("supervisor.actionBar")}
                </p>
                <div className="flex flex-wrap gap-2">
                  <button
                    type="button"
                    onClick={() => approveMut.mutate()}
                    disabled={busy}
                    className="font-amharic inline-flex items-center gap-1.5 rounded-md bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground hover:bg-primary/90 disabled:opacity-60"
                  >
                    {approveMut.isPending ? (
                      <Loader2 className="h-4 w-4 animate-spin" />
                    ) : (
                      <CheckCircle2 className="h-4 w-4" />
                    )}
                    {t("actions.finalApprove")}
                  </button>
                  <button
                    type="button"
                    onClick={() => setReturnOpen(true)}
                    disabled={busy}
                    className="font-amharic inline-flex items-center gap-1.5 rounded-md border border-destructive/40 bg-background px-4 py-2 text-sm font-medium text-destructive hover:bg-destructive/10 disabled:opacity-60"
                  >
                    <RotateCcw className="h-4 w-4" />
                    {t("actions.returnToSupervisor")}
                  </button>
                  {record.status !== "draft" && (
                    <button
                      type="button"
                      onClick={() => void handleDownloadPdf()}
                      disabled={pdfLoading || busy}
                      className="font-amharic inline-flex items-center gap-1.5 rounded-md border border-input bg-background px-4 py-2 text-sm font-medium text-foreground hover:bg-accent disabled:opacity-60"
                    >
                      {pdfLoading ? (
                        <Loader2 className="h-4 w-4 animate-spin" />
                      ) : (
                        <Download className="h-4 w-4" />
                      )}
                      {t("actions.downloadPdf")}
                    </button>
                  )}
                </div>
              </>
            )}
            {record.status !== "under_review" && record.status !== "draft" && (
              <button
                type="button"
                onClick={() => void handleDownloadPdf()}
                disabled={pdfLoading}
                className="font-amharic inline-flex items-center gap-1.5 rounded-md border border-input bg-background px-4 py-2 text-sm font-medium text-foreground hover:bg-accent disabled:opacity-60"
              >
                {pdfLoading ? (
                  <Loader2 className="h-4 w-4 animate-spin" />
                ) : (
                  <Download className="h-4 w-4" />
                )}
                {t("actions.downloadPdf")}
              </button>
            )}
          </div>
        )}
      </div>

      <RecordDetailTabs
        recordType={recordType}
        recordId={recordId}
        record={record}
        history={history}
        comments={comments}
        photos={photos}
        canComment
      />

      <ReturnModal
        open={returnOpen}
        onClose={() => setReturnOpen(false)}
        onConfirm={(comment) => returnMut.mutateAsync(comment)}
        title={t("actions.returnToSupervisor")}
        titleAm={t("actions.returnToSupervisor")}
      />
    </div>
  );
}
