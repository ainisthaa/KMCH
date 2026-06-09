# Convenience targets for resetting / inspecting the dev MySQL container.
# Run from this folder: `cd KMCH-BE/hospitalKMCH && make <target>`

MYSQL_CONTAINER ?= lineoa_mysql
MYSQL_USER      ?= root
MYSQL_PASS      ?= root
MYSQL_DB        ?= lineoa_miniapp
MYSQL          := docker exec $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER) -p$(MYSQL_PASS)

.PHONY: help db-reset db-wipe db-status db-shell db-delete db-mark-psyeval db-mark-sv up down logs

help:
	@echo "DB targets:"
	@echo "  make db-reset                       Delete all patient rows (keeps event & rooms seed)"
	@echo "  make db-wipe                        Full reset: drop docker volume + recreate (re-runs auto-migrate + seed)"
	@echo "  make db-status                      Row counts + 5 most recent patient_info"
	@echo "  make db-shell                       Open interactive mysql shell"
	@echo "  make db-delete LINE_ID=dev-...      Delete one patient by line_id"
	@echo "  make db-mark-psyeval LINE_ID=...    Set psyeval_form=1 (simulate completed mental form)"
	@echo "  make db-mark-sv LINE_ID=...         Set is_sv=1 (simulate Severe flag)"
	@echo "  make db-mark-paid LINE_ID=...       Set is_paid=1 (skip the payment QR scan)"
	@echo "  make db-join-queue LINE_ID=...      Mark paid + insert a waiting queue row (skip QR scan)"
	@echo ""
	@echo "App targets:"
	@echo "  make up                             docker compose up -d"
	@echo "  make down                           docker compose down (keeps volume)"
	@echo "  make logs                           tail backend logs"

db-reset:
	$(MYSQL) -e "USE $(MYSQL_DB); \
	  DELETE FROM patient_queue; \
	  DELETE FROM patient_check; \
	  DELETE FROM patient_info;"
	@echo "Patient tables cleared."

db-wipe:
	docker compose down -v
	docker compose up -d
	@echo "Volume destroyed; backend restarted (auto-migrate + seed will re-run)."

db-status:
	@$(MYSQL) -t -e "USE $(MYSQL_DB); \
	  SELECT (SELECT COUNT(*) FROM patient_info)  AS patient_info, \
	         (SELECT COUNT(*) FROM patient_check) AS patient_check, \
	         (SELECT COUNT(*) FROM patient_queue) AS patient_queue; \
	  SELECT line_id, first_name, last_name, register_date \
	    FROM patient_info ORDER BY register_date DESC LIMIT 5;"

db-shell:
	docker exec -it $(MYSQL_CONTAINER) mysql -u$(MYSQL_USER) -p$(MYSQL_PASS) $(MYSQL_DB)

db-delete:
	@test -n "$(LINE_ID)" || (echo "usage: make db-delete LINE_ID=dev-..." && exit 1)
	$(MYSQL) -e "USE $(MYSQL_DB); \
	  DELETE FROM patient_queue WHERE line_id='$(LINE_ID)'; \
	  DELETE FROM patient_check WHERE line_id='$(LINE_ID)'; \
	  DELETE FROM patient_info  WHERE line_id='$(LINE_ID)';"
	@echo "Deleted $(LINE_ID)."

db-mark-psyeval:
	@test -n "$(LINE_ID)" || (echo "usage: make db-mark-psyeval LINE_ID=..." && exit 1)
	$(MYSQL) -e "USE $(MYSQL_DB); UPDATE patient_check SET psyeval_form=1 WHERE line_id='$(LINE_ID)';"
	@echo "psyeval_form=1 for $(LINE_ID)."

db-mark-sv:
	@test -n "$(LINE_ID)" || (echo "usage: make db-mark-sv LINE_ID=..." && exit 1)
	$(MYSQL) -e "USE $(MYSQL_DB); UPDATE patient_check SET is_sv=1 WHERE line_id='$(LINE_ID)';"
	@echo "is_sv=1 for $(LINE_ID)."

db-mark-paid:
	@test -n "$(LINE_ID)" || (echo "usage: make db-mark-paid LINE_ID=..." && exit 1)
	$(MYSQL) -e "USE $(MYSQL_DB); UPDATE patient_check SET is_paid=1 WHERE line_id='$(LINE_ID)';"
	@echo "is_paid=1 for $(LINE_ID)."

# Force a registered+paid patient into the doctor-consultation queue without
# needing to scan QR codes from the phone. Useful for testing the Queue page.
db-join-queue:
	@test -n "$(LINE_ID)" || (echo "usage: make db-join-queue LINE_ID=..." && exit 1)
	$(MYSQL) -e "USE $(MYSQL_DB); \
	  UPDATE patient_check SET is_paid=1 WHERE line_id='$(LINE_ID)'; \
	  INSERT INTO patient_queue (line_id, event_id, queue, status, station, q_starttime, created_at, updated_at) \
	    VALUES ('$(LINE_ID)', 1, DATE_FORMAT(NOW(),'%Y%m%d%H%i%s'), 'waiting', 'queue', NOW(), NOW(), NOW());"
	@echo "$(LINE_ID) added to queue (status=waiting). Auto-fill will assign a room if free."

up:
	docker compose up -d

down:
	docker compose down

logs:
	docker compose logs -f app
