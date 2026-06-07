import { useEffect, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { isAxiosError } from "axios";
import { createUser, updateUser } from "@/api/users";
import type { CreateUserBody, Language, Role, UpdateUserBody, UserItem } from "@/types";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Field, Select, Switch, TextInput } from "@/components/forms/Field";

interface UserFormModalProps {
  open: boolean;
  onClose: () => void;
  user?: UserItem;
}

export function UserFormModal({ open, onClose, user }: UserFormModalProps) {
  const { t } = useTranslation();
  const qc = useQueryClient();
  const isEdit = Boolean(user);

  const [fullName, setFullName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [role, setRole] = useState<Role>("registrar");
  const [language, setLanguage] = useState<Language>("am");
  const [isActive, setIsActive] = useState(true);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [apiError, setApiError] = useState("");

  useEffect(() => {
    if (!open) return;
    setErrors({});
    setApiError("");
    if (user) {
      setFullName(user.full_name);
      setEmail(user.email);
      setPassword("");
      setRole(user.role);
      setLanguage(user.language);
      setIsActive(user.is_active);
    } else {
      setFullName("");
      setEmail("");
      setPassword("");
      setRole("registrar");
      setLanguage("am");
      setIsActive(true);
    }
  }, [open, user]);

  const validate = () => {
    const next: Record<string, string> = {};
    if (!fullName.trim()) next.fullName = t("users.requiredName");
    if (fullName.length > 100) next.fullName = t("users.nameTooLong");
    if (!isEdit) {
      if (!email.trim() || !email.includes("@")) next.email = t("users.invalidEmail");
      if (password.length < 8) next.password = t("users.passwordTooShort");
    }
    if (!role) next.role = t("users.requiredRole");
    setErrors(next);
    return Object.keys(next).length === 0;
  };

  const saveMut = useMutation({
    mutationFn: async () => {
      if (isEdit && user) {
        const body: UpdateUserBody = {
          full_name: fullName.trim(),
          role,
          language,
          is_active: isActive,
        };
        return updateUser(user.id, body);
      }
      const body: CreateUserBody = {
        full_name: fullName.trim(),
        email: email.trim(),
        password,
        role,
        language,
      };
      return createUser(body);
    },
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["users"] });
      toast.success(isEdit ? t("toast.updateUserSuccess") : t("toast.createUserSuccess"));
      onClose();
    },
    onError: (err) => {
      if (isAxiosError(err) && err.response?.status === 409) {
        setApiError(t("users.emailExists"));
        return;
      }
      toast.error(t("toast.error"));
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setApiError("");
    if (!validate()) return;
    saveMut.mutate();
  };

  const roleOptions = [
    { value: "registrar", labelAm: t("roles.registrar"), labelEn: "Registrar" },
    { value: "supervisor", labelAm: t("roles.supervisor"), labelEn: "Supervisor" },
    { value: "manager", labelAm: t("roles.manager"), labelEn: "Manager" },
  ];

  const langOptions = [
    { value: "am", labelAm: t("common.amharic") },
    { value: "en", labelAm: t("common.english") },
  ];

  return (
    <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
      <DialogContent className="max-h-[90vh] overflow-y-auto sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="font-amharic">
            {isEdit ? t("users.editTitle") : t("users.createTitle")}
          </DialogTitle>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          <Field am={t("users.fullName")} en="Full Name" required>
            <TextInput
              value={fullName}
              onChange={(e) => setFullName(e.target.value)}
              error={errors.fullName}
            />
          </Field>

          {!isEdit && (
            <Field am={t("users.email")} en="Email" required>
              <TextInput
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                error={errors.email}
              />
            </Field>
          )}

          {isEdit && (
            <Field am={t("users.email")} en="Email">
              <TextInput value={email} readOnly disabled className="opacity-70" />
            </Field>
          )}

          {!isEdit && (
            <Field am={t("users.password")} en="Password" required>
              <TextInput
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                error={errors.password}
              />
            </Field>
          )}

          <Field am={t("users.role")} en="Role" required>
            <Select
              value={role}
              onChange={(e) => setRole(e.target.value as Role)}
              options={roleOptions}
              error={errors.role}
            />
          </Field>

          <Field am={t("users.language")} en="Language">
            <Select
              value={language}
              onChange={(e) => setLanguage(e.target.value as Language)}
              options={langOptions}
            />
          </Field>

          {isEdit && (
            <Switch
              checked={isActive}
              onChange={setIsActive}
              am={t("users.isActive")}
              en="Active"
            />
          )}

          {apiError && (
            <p className="text-sm text-destructive">{apiError}</p>
          )}

          <DialogFooter className="gap-2 sm:gap-0">
            <button
              type="button"
              onClick={onClose}
              disabled={saveMut.isPending}
              className="rounded-md border border-input bg-background px-4 py-2 text-sm font-medium hover:bg-accent disabled:opacity-60"
            >
              {t("common.cancel")}
            </button>
            <button
              type="submit"
              disabled={saveMut.isPending}
              className="inline-flex items-center gap-1.5 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-60"
            >
              {saveMut.isPending && <Loader2 className="h-4 w-4 animate-spin" />}
              {isEdit ? t("common.save") : t("users.newUser")}
            </button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
