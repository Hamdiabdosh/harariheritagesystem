import { useState } from "react";
import { useQuery, useMutation, useQueryClient, keepPreviousData } from "@tanstack/react-query";
import { useTranslation } from "react-i18next";
import { Plus, Pencil, UserX, ChevronLeft, ChevronRight } from "lucide-react";
import { toast } from "sonner";
import { listUsers, deactivateUser } from "@/api/users";
import { useAuthStore } from "@/stores/authStore";
import type { Role, UserItem, UserListParams } from "@/types";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { EmptyState } from "@/components/common/EmptyState";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { UserFormModal } from "./UserFormModal";
import { ConfirmDeactivateModal } from "./ConfirmDeactivateModal";

const ROLE_BADGE: Record<Role, string> = {
  registrar: "bg-blue-100 text-blue-900",
  supervisor: "bg-amber-100 text-amber-900",
  manager: "bg-purple-100 text-purple-900",
};

export function UserManagement() {
  const { t } = useTranslation();
  const qc = useQueryClient();
  const currentUser = useAuthStore((s) => s.user);

  const [filters, setFilters] = useState<UserListParams>({ page: 1, limit: 20 });
  const [formOpen, setFormOpen] = useState(false);
  const [editUser, setEditUser] = useState<UserItem | undefined>();
  const [deactivateTarget, setDeactivateTarget] = useState<UserItem | null>(null);

  const { data, isLoading, isError, error, refetch } = useQuery({
    queryKey: ["users", filters],
    queryFn: () => listUsers(filters),
    placeholderData: keepPreviousData,
  });

  const deactivateMut = useMutation({
    mutationFn: (id: string) => deactivateUser(id),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["users"] });
      toast.success(t("toast.deactivateSuccess"));
      setDeactivateTarget(null);
    },
    onError: () => toast.error(t("toast.error")),
  });

  const openCreate = () => {
    setEditUser(undefined);
    setFormOpen(true);
  };

  const openEdit = (user: UserItem) => {
    setEditUser(user);
    setFormOpen(true);
  };

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex flex-wrap items-center gap-2">
          <select
            value={filters.role ?? ""}
            onChange={(e) =>
              setFilters((f) => ({
                ...f,
                page: 1,
                role: (e.target.value || undefined) as Role | undefined,
              }))
            }
            className="rounded-md border border-input bg-background px-3 py-2 text-sm"
          >
            <option value="">{t("users.filterAll")}</option>
            <option value="registrar">{t("roles.registrar")}</option>
            <option value="supervisor">{t("roles.supervisor")}</option>
            <option value="manager">{t("roles.manager")}</option>
          </select>

          <select
            value={
              filters.is_active === undefined ? "" : filters.is_active ? "active" : "inactive"
            }
            onChange={(e) => {
              const v = e.target.value;
              setFilters((f) => ({
                ...f,
                page: 1,
                is_active: v === "" ? undefined : v === "active",
              }));
            }}
            className="rounded-md border border-input bg-background px-3 py-2 text-sm"
          >
            <option value="">{t("users.allUsers")}</option>
            <option value="active">{t("users.active")}</option>
            <option value="inactive">{t("users.inactive")}</option>
          </select>
        </div>

        <button
          type="button"
          onClick={openCreate}
          className="font-amharic inline-flex items-center gap-1.5 rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90"
        >
          <Plus className="h-4 w-4" />
          {t("users.newUser")}
        </button>
      </div>

      {isLoading && (
        <div className="flex justify-center rounded-xl border border-border bg-card p-12">
          <LoadingSpinner />
        </div>
      )}

      {isError && (
        <div className="rounded-xl border border-rose-200 bg-rose-50 p-4 text-sm text-rose-900">
          {(error as Error)?.message}
          <button onClick={() => refetch()} className="ml-2 underline">
            {t("common.retry")}
          </button>
        </div>
      )}

      {data && data.items.length === 0 && (
        <EmptyState title={t("common.empty")} description={t("users.allUsers")} />
      )}

      {data && data.items.length > 0 && (
        <>
          <div className="overflow-hidden rounded-xl border border-border bg-card">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="font-amharic">{t("users.fullName")}</TableHead>
                  <TableHead>{t("users.email")}</TableHead>
                  <TableHead className="font-amharic">{t("users.role")}</TableHead>
                  <TableHead className="font-amharic">{t("users.isActive")}</TableHead>
                  <TableHead className="text-right">{t("users.actions")}</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.items.map((user) => {
                  const isSelf = user.id === currentUser?.id;
                  return (
                    <TableRow key={user.id}>
                      <TableCell className="font-amharic font-medium">{user.full_name}</TableCell>
                      <TableCell className="text-sm text-muted-foreground">{user.email}</TableCell>
                      <TableCell>
                        <span
                          className={`font-amharic inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${ROLE_BADGE[user.role]}`}
                        >
                          {t(`roles.${user.role}`)}
                        </span>
                      </TableCell>
                      <TableCell>
                        <span
                          className={`font-amharic inline-flex rounded-full px-2 py-0.5 text-xs font-medium ${
                            user.is_active
                              ? "bg-emerald-100 text-emerald-900"
                              : "bg-muted text-muted-foreground"
                          }`}
                        >
                          {user.is_active ? t("users.active") : t("users.inactive")}
                        </span>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="inline-flex gap-1">
                          <button
                            type="button"
                            onClick={() => openEdit(user)}
                            className="inline-flex h-8 w-8 items-center justify-center rounded-md border border-input hover:bg-accent"
                            aria-label={t("users.editTitle")}
                          >
                            <Pencil className="h-3.5 w-3.5" />
                          </button>
                          <button
                            type="button"
                            onClick={() => setDeactivateTarget(user)}
                            disabled={isSelf || !user.is_active}
                            className="inline-flex h-8 w-8 items-center justify-center rounded-md border border-destructive/30 text-destructive hover:bg-destructive/10 disabled:opacity-40"
                            aria-label={t("users.confirmDeactivate")}
                          >
                            <UserX className="h-3.5 w-3.5" />
                          </button>
                        </div>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </div>

          <div className="flex flex-wrap items-center justify-between gap-2">
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
                onClick={() => setFilters((f) => ({ ...f, page: (f.page ?? 1) - 1 }))}
                className="inline-flex h-8 items-center gap-1 rounded-md border border-input px-2 text-xs disabled:opacity-40"
              >
                <ChevronLeft className="h-3.5 w-3.5" />
                {t("common.prev")}
              </button>
              <button
                type="button"
                disabled={data.page >= data.total_pages}
                onClick={() => setFilters((f) => ({ ...f, page: (f.page ?? 1) + 1 }))}
                className="inline-flex h-8 items-center gap-1 rounded-md border border-input px-2 text-xs disabled:opacity-40"
              >
                {t("common.next")}
                <ChevronRight className="h-3.5 w-3.5" />
              </button>
            </div>
          </div>
        </>
      )}

      <UserFormModal open={formOpen} onClose={() => setFormOpen(false)} user={editUser} />

      {deactivateTarget && (
        <ConfirmDeactivateModal
          open
          onClose={() => setDeactivateTarget(null)}
          userName={deactivateTarget.full_name}
          onConfirm={() => deactivateMut.mutateAsync(deactivateTarget.id)}
        />
      )}
    </div>
  );
}
