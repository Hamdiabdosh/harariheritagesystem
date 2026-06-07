import { useState } from "react";
import { useQuery, keepPreviousData } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { FileText, ChevronLeft, ChevronRight } from "lucide-react";
import type { UnifiedListParams, PaginatedRecordSummaries } from "@/types";
import { RecordCard } from "./RecordCard";
import { RecordFilters } from "./RecordFilters";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { EmptyState } from "@/components/common/EmptyState";

interface RecordsListProps {
  queryKey: readonly unknown[];
  fetcher: (params: UnifiedListParams) => Promise<PaginatedRecordSummaries>;
  initialParams?: UnifiedListParams;
  detailPath?: string;
}

export function RecordsList({
  queryKey,
  fetcher,
  initialParams = {},
  detailPath,
}: RecordsListProps) {
  const { t } = useTranslation();
  const [params, setParams] = useState<UnifiedListParams>({
    page: 1,
    limit: 20,
    ...initialParams,
  });

  const { data, isLoading, isFetching, isError, error, refetch } = useQuery({
    queryKey: [...queryKey, params],
    queryFn: () => fetcher(params),
    placeholderData: keepPreviousData,
  });

  return (
    <div className="space-y-4">
      <RecordFilters value={params} onChange={setParams} />

      {isLoading && (
        <div className="flex justify-center rounded-xl border border-border bg-card p-12">
          <LoadingSpinner />
        </div>
      )}

      {isError && (
        <div className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-900">
          <div className="font-medium">{t("records.loadError")}</div>
          <div className="mt-1 text-xs opacity-80">
            {(error as Error)?.message}
          </div>
          <button
            onClick={() => refetch()}
            className="mt-3 inline-flex items-center rounded-md border border-rose-300 bg-white px-3 py-1.5 text-xs font-medium hover:bg-rose-100"
          >
            {t("common.retry")}
          </button>
        </div>
      )}

      {data && data.items.length === 0 && (
        <EmptyState
          icon={<FileText className="h-6 w-6" />}
          title={t("records.emptyTitle")}
          description={t("records.emptyDescription")}
        />
      )}

      {data && data.items.length > 0 && (
        <>
          <div
            className={`grid gap-3 ${isFetching ? "opacity-60 transition-opacity" : ""}`}
          >
            {data.items.map((rec) => (
              <RecordCard key={rec.id} record={rec} detailPath={detailPath} />
            ))}
          </div>

          {/* Pagination */}
          <div className="flex flex-wrap items-center justify-between gap-2 pt-2">
            <div className="text-xs text-muted-foreground tabular-nums">
              {t("records.pagination", {
                shown: data.items.length,
                total: data.total,
                page: data.page,
                pages: data.total_pages,
              })}
            </div>
            <div className="flex items-center gap-1">
              <button
                type="button"
                disabled={data.page <= 1}
                onClick={() =>
                  setParams((p) => ({ ...p, page: (p.page ?? 1) - 1 }))
                }
                className="inline-flex h-8 items-center gap-1 rounded-md border border-input bg-background px-2 text-xs text-foreground transition-colors hover:bg-accent disabled:opacity-40"
              >
                <ChevronLeft className="h-3.5 w-3.5" />
                {t("common.prev")}
              </button>
              <button
                type="button"
                disabled={data.page >= data.total_pages}
                onClick={() =>
                  setParams((p) => ({ ...p, page: (p.page ?? 1) + 1 }))
                }
                className="inline-flex h-8 items-center gap-1 rounded-md border border-input bg-background px-2 text-xs text-foreground transition-colors hover:bg-accent disabled:opacity-40"
              >
                {t("common.next")}
                <ChevronRight className="h-3.5 w-3.5" />
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
