import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { MessageSquare, Loader2 } from "lucide-react";
import type { RecordComment, RecordType } from "@/types";
import { addComment } from "@/api/workflow";
import { EmptyState } from "@/components/common/EmptyState";

interface CommentThreadProps {
  comments: RecordComment[];
  recordType: RecordType;
  recordId: string;
  canComment: boolean;
  canViewComments?: boolean;
  readOnly?: boolean;
}

function formatDate(iso: string, locale: string) {
  try {
    return new Date(iso).toLocaleDateString(locale === "am" ? "am-ET" : "en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  } catch {
    return iso;
  }
}

export function CommentThread({
  comments,
  recordType,
  recordId,
  canComment,
  canViewComments = true,
  readOnly = false,
}: CommentThreadProps) {
  const { t, i18n } = useTranslation();
  const qc = useQueryClient();
  const [text, setText] = useState("");
  const [error, setError] = useState<string | null>(null);

  const addMut = useMutation({
    mutationFn: (commentText: string) => addComment(recordType, recordId, commentText),
    onSuccess: () => {
      setText("");
      setError(null);
      qc.invalidateQueries({ queryKey: ["record-detail", recordType, recordId] });
      qc.invalidateQueries({ queryKey: [recordType, recordId] });
    },
    onError: () => setError(t("comments.error")),
  });

  const sorted = [...comments].sort(
    (a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime(),
  );

  const handleSubmit = () => {
    const trimmed = text.trim();
    if (!trimmed) return;
    addMut.mutate(trimmed);
  };

  if (comments.length === 0 && readOnly) {
    return (
      <EmptyState
        icon={<MessageSquare className="h-6 w-6" />}
        title={t("comments.empty")}
      />
    );
  }

  return (
    <div className="space-y-4">
      {comments.length === 0 && canComment && !readOnly && (
        <p className="font-amharic text-sm text-muted-foreground">{t("comments.add")}</p>
      )}

      {canViewComments && sorted.length > 0 && (
        <ul className="space-y-3">
          {sorted.map((c) => (
            <li
              key={c.id}
              className="rounded-lg border border-border bg-card p-3"
            >
              <div className="flex items-baseline justify-between gap-2">
                <span className="text-xs font-medium text-muted-foreground">
                  {c.author_name}
                </span>
                <time className="text-xs text-muted-foreground">
                  {formatDate(c.created_at, i18n.language)}
                </time>
              </div>
              <p className="mt-1.5 text-sm text-foreground">{c.comment_text}</p>
            </li>
          ))}
        </ul>
      )}

      {canComment && !readOnly && (
        <div className="space-y-2 border-t border-border pt-4">
          <textarea
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder={t("comments.placeholder")}
            rows={3}
            className="font-amharic w-full rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
          {error && <p className="text-xs text-destructive">{error}</p>}
          <button
            type="button"
            onClick={handleSubmit}
            disabled={addMut.isPending || !text.trim()}
            className="font-amharic inline-flex items-center gap-1.5 rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-60"
          >
            {addMut.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
            {t("comments.submit")}
          </button>
        </div>
      )}
    </div>
  );
}
