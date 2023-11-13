package ar.com.montepagano.search.examples.creditcardpayments;

import java.io.FileInputStream;
import java.io.IOException;
import java.util.HashMap;
import java.util.Map;
import java.util.Scanner;

import com.google.gson.Gson;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import com.google.protobuf.ByteString;
import ar.com.montepagano.search.v1.PrivateMiddlewareServiceGrpc;
import ar.com.montepagano.search.v1.Middleware;
import ar.com.montepagano.search.v1.Middleware.RegisterChannelRequest;
import ar.com.montepagano.search.v1.Middleware.RegisterChannelResponse;
import ar.com.montepagano.search.v1.AppMessageOuterClass.AppSendRequest;
import ar.com.montepagano.search.v1.AppMessageOuterClass.AppMessage;
import ar.com.montepagano.search.v1.Contracts.GlobalContract;
import ar.com.montepagano.search.v1.Contracts.GlobalContractFormat;
import ar.com.montepagano.search.v1.Broker.RemoteParticipant;

public class Main {
    public static void main(String[] args) {
        Scanner scanner = new Scanner(System.in);

        // List of book titles
        String[] bookTitles = {
                "The Great Gatsby",
                "To Kill a Mockingbird",
                "1984",
                "Pride and Prejudice",
                "Animal Farm",
                "The Catcher in the Rye",
                "Lord of the Flies",
                "Brave New World",
                "The Hobbit",
                "The Fellowship of the Ring"
        };

        // Prompt user to select book titles
        System.out.println("Please select at least one book to purchase:");
        for (int i = 0; i < bookTitles.length; i++) {
            System.out.println((i + 1) + ". " + bookTitles[i]);
        }

        // Get user input for selected book titles
        System.out.print("Enter the number(s) of the book(s) you want to purchase (comma-separated): ");
        String selectedBooksInput = scanner.nextLine();
        String[] selectedBooksInputArray = selectedBooksInput.split(",");
        int[] selectedBooks = new int[selectedBooksInputArray.length];
        for (int i = 0; i < selectedBooksInputArray.length; i++) {
            selectedBooks[i] = Integer.parseInt(selectedBooksInputArray[i].trim());
        }

        // Prompt user for shipping address
        System.out.print("Enter your shipping address: ");
        String shippingAddress = scanner.nextLine();

        // Print purchase information
        System.out.println("Purchase Information:");
        System.out.println("Selected Books:");
        for (int selectedBook : selectedBooks) {
            System.out.println("- " + bookTitles[selectedBook - 1]);
        }
        System.out.println("Shipping Address: " + shippingAddress);

        // Load file contract.fsa into a GlobalContract
        ByteString contractBytes = null;
        try {
            contractBytes = ByteString.readFrom(new FileInputStream("contract.fsa"));
        } catch (IOException e) {
            e.printStackTrace();
        }
        GlobalContract contract = GlobalContract.newBuilder().setContract(contractBytes).setFormat(
                GlobalContractFormat.GLOBAL_CONTRACT_FORMAT_FSA).setInitiatorName("ClientApp").build();

        // Get the stub to communicate with the middleware
        ManagedChannel channel = ManagedChannelBuilder.forTarget(
                "middleware-client:11000").usePlaintext().build();
        PrivateMiddlewareServiceGrpc.PrivateMiddlewareServiceBlockingStub stub =
                PrivateMiddlewareServiceGrpc.newBlockingStub(channel);

        // Prompt the user to input hostname for the PPS and Srv apps. Both should have a default value of "middleware-payments:10000" and "middleware-backend:10000" respectively.
        System.out.print("Enter the hostname for the PPS app (default: middleware-payments:10000): ");
        String ppsHostname = scanner.nextLine();
        if (ppsHostname.isEmpty()) {
            ppsHostname = "middleware-payments:10000";
        }
        System.out.print("Enter the hostname for the Srv app (default: middleware-backend:10000): ");
        String srvHostname = scanner.nextLine();
        if (srvHostname.isEmpty()) {
            srvHostname = "middleware-backend:10000";
        }
        // Prompt the user to input AppId for the PPS and Srv apps. They don't have default values and must be provided by the user.
        System.out.print("Enter the AppId for the PPS app: ");
        String ppsAppId = scanner.nextLine();
        System.out.print("Enter the AppId for the Srv app: ");
        String srvAppId = scanner.nextLine();

        // Register channel with middleware using GlobalContract.
        // Add preset participants to the RegisterChannelRequest. We do this because we don't have an algorithm for compatibility checking in the broker.
        RemoteParticipant pps = RemoteParticipant.newBuilder().setUrl(ppsHostname).setAppId(ppsAppId).build();
        RemoteParticipant srv = RemoteParticipant.newBuilder().setUrl(srvHostname).setAppId(srvAppId).build();
        RegisterChannelRequest request = RegisterChannelRequest.newBuilder().setRequirementsContract(contract).putPresetParticipants("PPS", pps).putPresetParticipants("Srv", srv).build();
        RegisterChannelResponse response = stub.registerChannel(request);
        var channelId = response.getChannelId();

        // Send PurchaseRequest with each item quantities and the shipping address
        Map<String, Integer> items = new HashMap<>();
        for (int selectedBook : selectedBooks) {
            String bookTitle = bookTitles[selectedBook - 1];
            if (items.containsKey(bookTitle)) {
                items.put(bookTitle, items.get(bookTitle) + 1);
            } else {
                items.put(bookTitle, 1);
            }
        }
        Gson gson = new Gson();
        var body = ByteString.copyFromUtf8(String.format(
                "{\"items\": %s, \"shippingAddress\": \"%s\"}",
                gson.toJson(items), shippingAddress));
        var msg = AppMessage.newBuilder().setType("PurchaseRequest").setBody(body).build();
        var sendreq = AppSendRequest.newBuilder().setChannelId(channelId).setRecipient("Srv").setMessage(msg).build();
        var sendresp = stub.appSend(sendreq);
        if (sendresp.getResult() != Middleware.AppSendResponse.Result.RESULT_OK) {
            System.out.println("Error sending PurchaseRequest. Exiting...");
            System.exit(1);
        }

        // Receive TotalAmount from Srv
        var recvreq = Middleware.AppRecvRequest.newBuilder().setChannelId(channelId).setParticipant("Srv").build();
        var recvresp = stub.appRecv(recvreq);
        if (!recvresp.getMessage().getType().equals("TotalAmount")) {
            System.out.println("Error receiving TotalAmount. Exiting...");
            System.exit(1);
        }
        var total_amount = gson.fromJson(recvresp.getMessage().getBody().toStringUtf8(), double.class);
        System.out.println("Total amount: " + total_amount);

        // Ask the user for the credit card details.
        System.out.print("Enter the credit card number: ");
        String cardNumber = scanner.nextLine();
        System.out.print("Enter the credit card expiration date: ");
        String cardExpirationDate = scanner.nextLine();
        System.out.print("Enter the credit card security code: ");
        String cardSecurityCode = scanner.nextLine();

        // Send CardDetailsWithTotalAmount to PPS.
        body = ByteString.copyFromUtf8(String.format("{\"card_number\": \"%s\", \"card_expirationDate\": \"%s\", \"card_cvv\": \"%s\", \"total_amount\": %s}", cardNumber, cardExpirationDate, cardSecurityCode, total_amount));
        msg = AppMessage.newBuilder().setType("CardDetailsWithTotalAmount").setBody(body).build();
        sendreq = AppSendRequest.newBuilder().setChannelId(channelId).setRecipient("PPS").setMessage(msg).build();
        sendresp = stub.appSend(sendreq);
        if (sendresp.getResult() != Middleware.AppSendResponse.Result.RESULT_OK) {
            System.out.println("Error sending CardDetailsWithTotalAmount. Exiting...");
            System.exit(1);
        }

        // Receive PaymentNonce from PPS.
        recvreq = Middleware.AppRecvRequest.newBuilder().setChannelId(channelId).setParticipant("PPS").build();
        recvresp = stub.appRecv(recvreq);
        if (!recvresp.getMessage().getType().equals("PaymentNonce")) {
            System.out.println("Error receiving PaymentNonce. Exiting...");
            System.exit(1);
        }
        var payment_nonce = gson.fromJson(recvresp.getMessage().getBody().toStringUtf8(), String.class);
        System.out.println("Payment nonce: " + payment_nonce);

        // Send PurchaseWithPaymentNonce to Srv.
        body = ByteString.copyFromUtf8(String.format("{\"items\": %s, \"shippingAddress\": \"%s\", \"nonce\": \"%s\"}", gson.toJson(items), shippingAddress, payment_nonce));
        msg = AppMessage.newBuilder().setType("PurchaseWithPaymentNonce").setBody(body).build();
        sendreq = AppSendRequest.newBuilder().setChannelId(channelId).setRecipient("Srv").setMessage(msg).build();
        sendresp = stub.appSend(sendreq);
        if (sendresp.getResult() != Middleware.AppSendResponse.Result.RESULT_OK) {
            System.out.println("Error sending PurchaseWithPaymentNonce. Exiting...");
            System.exit(1);
        }

        // Receive either PurchaseOK or PurchaseFail from Srv
        recvreq = Middleware.AppRecvRequest.newBuilder().setChannelId(channelId).setParticipant("Srv").build();
        recvresp = stub.appRecv(recvreq);
        if (recvresp.getMessage().getType().equals("PurchaseOK")) {
            System.out.println("Purchase successful!");
        } else if (recvresp.getMessage().getType().equals("PurchaseFail")) {
            System.out.println("Purchase failed!");
        } else {
            System.out.println("Error receiving PurchaseOK or PurchaseFail. Exiting...");
            System.exit(1);
        }

        channel.shutdown();
    }
}