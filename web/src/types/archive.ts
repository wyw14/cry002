export type CaseStatus='draft'|'pending_review'|'archived'|'borrowed'|'sealed'|'destroyed'
export interface CaseFile{id:string;case_number:string;title:string;classification_id:string;security_level:string;retention:string;responsible_department:string;handler_id:string;year:number;keywords:string[];summary:string;status:CaseStatus;version:number;updated_at:string}
export interface BorrowRequest{id:string;case_id:string;applicant_id:string;purpose:string;status:string;due_at:string;renewals:number}
export interface Classification{id:string;parent_id?:string;name:string;code:string;sort_order:number;enabled:boolean}
export interface AuditEvent{id:string;actor_id:string;action:string;resource:string;resource_id:string;request_id:string;created_at:string}
export interface Page<T>{items:T[];total:number;page:number;page_size:number}
