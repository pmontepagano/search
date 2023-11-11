
# Use Python 3.12 as base image
FROM python:3.12
ENV PYTHONUNBUFFERED=1

# Set working directory to /app
WORKDIR /app

# WARNING! Paths in the Dockerfile are relative to the root of the repository.
# We do that to be able to copy the generated stubs from the gen directory.
# The context is set in the docker-compose.yml file.
COPY examples/credit-card-payments/backend/requirements.txt .
# Install dependencies from requirements.txt
RUN pip install -r requirements.txt

COPY examples/credit-card-payments/backend/ .
RUN rm lib
COPY gen/python lib



# Set main.py as default command
CMD ["python", "main.py"]
