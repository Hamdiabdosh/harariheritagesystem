import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Loader2 } from "lucide-react";
import { useTranslation } from "react-i18next";

interface ConfirmDeactivateModalProps {
  open: boolean;
  onClose: () => void;
  onConfirm: () => Promise<void>;
  userName: string;
}

export function ConfirmDeactivateModal({
  open,
  onClose,
  onConfirm,
  userName,
}: ConfirmDeactivateModalProps) {
  const { t } = useTranslation();
  const [pending, setPending] = useState(false);

  const handleConfirm = async () => {
    setPending(true);
    try {
      await onConfirm();
      onClose();
    } finally {
      setPending(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle className="font-amharic">{t("users.confirmDeactivate")}</DialogTitle>
        </DialogHeader>
        <p className="font-amharic text-sm text-muted-foreground">
          {t("users.deactivateBody", { name: userName })}
        </p>
        <p className="text-sm text-muted-foreground">{t("users.deactivateWarning")}</p>
        <DialogFooter className="gap-2 sm:gap-0">
          <button
            type="button"
            onClick={onClose}
            disabled={pending}
            className="rounded-md border border-input bg-background px-4 py-2 text-sm font-medium hover:bg-accent disabled:opacity-60"
          >
            {t("common.cancel")}
          </button>
          <button
            type="button"
            onClick={() => void handleConfirm()}
            disabled={pending}
            className="inline-flex items-center gap-1.5 rounded-md bg-destructive px-4 py-2 text-sm font-medium text-destructive-foreground hover:bg-destructive/90 disabled:opacity-60"
          >
            {pending && <Loader2 className="h-4 w-4 animate-spin" />}
            {t("users.confirmDeactivate")}
          </button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
