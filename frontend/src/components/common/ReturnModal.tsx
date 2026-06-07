import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Loader2 } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface ReturnModalProps {
  open: boolean;
  onClose: () => void;
  onConfirm: (comment: string) => Promise<void>;
  title: string;
  titleAm: string;
}

export function ReturnModal({
  open,
  onClose,
  onConfirm,
  title,
  titleAm,
}: ReturnModalProps) {
  const { t, i18n } = useTranslation();
  const [comment, setComment] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [pending, setPending] = useState(false);

  const handleClose = () => {
    if (pending) return;
    setComment("");
    setError(null);
    onClose();
  };

  const handleConfirm = async () => {
    if (!comment.trim()) {
      setError(t("modal.reasonRequired"));
      return;
    }
    setPending(true);
    setError(null);
    try {
      await onConfirm(comment.trim());
      setComment("");
      onClose();
    } catch {
      setError(t("modal.failed"));
    } finally {
      setPending(false);
    }
  };

  const isAm = i18n.language === "am";

  return (
    <Dialog open={open} onOpenChange={(v) => !v && handleClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="font-amharic">
            {isAm ? titleAm : title}
          </DialogTitle>
        </DialogHeader>
        <div className="space-y-2">
          <label className="font-amharic block text-sm font-medium text-foreground">
            {isAm ? t("modal.reasonLabelAm") : t("modal.reasonLabel")}
          </label>
          <textarea
            value={comment}
            onChange={(e) => {
              setComment(e.target.value);
              if (error) setError(null);
            }}
            rows={4}
            className="font-amharic w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          />
          {error && <p className="text-xs text-destructive">{error}</p>}
        </div>
        <DialogFooter>
          <button
            type="button"
            onClick={handleClose}
            disabled={pending}
            className="font-amharic rounded-md border border-input px-3 py-2 text-sm hover:bg-accent disabled:opacity-60"
          >
            {t("modal.cancel")}
          </button>
          <button
            type="button"
            onClick={handleConfirm}
            disabled={pending}
            className="font-amharic inline-flex items-center gap-1.5 rounded-md bg-destructive px-3 py-2 text-sm font-medium text-destructive-foreground hover:bg-destructive/90 disabled:opacity-60"
          >
            {pending && <Loader2 className="h-4 w-4 animate-spin" />}
            {t("modal.confirm")}
          </button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
