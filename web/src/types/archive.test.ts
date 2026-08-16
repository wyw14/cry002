import{describe,expect,it}from'vitest';import type{CaseStatus}from'./archive'
describe('case lifecycle',()=>{it('keeps the configured state order',()=>{const states:CaseStatus[]=['draft','pending_review','archived','borrowed','sealed','destroyed'];expect(states).toHaveLength(6);expect(states.indexOf('archived')).toBeLessThan(states.indexOf('sealed'))})})
